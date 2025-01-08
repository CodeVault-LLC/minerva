package core

import (
	"strings"
	"sync"
	"time"

	"github.com/codevault-llc/minerva/internal/common"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type PageAnalysis struct {
	browser *rod.Browser
	page    *rod.Page
}

func NewPageAnalysis() *PageAnalysis {
	browser := rod.New().MustConnect().NoDefaultDevice().MustIncognito()
	page := browser.MustPage()

	return &PageAnalysis{
		browser: browser,
		page:    page,
	}
}

func (p *PageAnalysis) Close() {
	p.page.Close()
	p.browser.MustClose()
}

type WebsiteResponse struct {
	Redirects    []common.Redirect
	Files        []common.FileRequest
	FinalHTML    string
	ParsedHTML   *html.Node
	WebsiteTitle string
	StatusCode   int
}

type TrackerQueue struct {
	src      string
	fileType utils.FileType
	headers  string
	cookies  string
	duration int
	status   int
}

func (p *PageAnalysis) FetchWebsite(url, userAgent string) (*WebsiteResponse, error) {
	requestTracker := make(map[proto.NetworkRequestID]TrackerQueue)
	var requestTrackerMutex sync.Mutex
	var redirects []common.Redirect
	var networkFiles []common.FileRequest

	addRequestID := func(id proto.NetworkRequestID, src string, fileType utils.FileType, clientHeaders string) {
		requestTrackerMutex.Lock()
		defer requestTrackerMutex.Unlock()
		requestTracker[id] = TrackerQueue{
			src:      src,
			fileType: fileType,
			headers:  clientHeaders,
		}
	}

	editStatus := func(id proto.NetworkRequestID, status int) {
		requestTrackerMutex.Lock()
		defer requestTrackerMutex.Unlock()
		queue, exists := requestTracker[id]
		if exists {
			queue.status = status
			requestTracker[id] = queue
		}
	}

	editRequest := func(id proto.NetworkRequestID, headers string, cookies string) {
		requestTrackerMutex.Lock()
		defer requestTrackerMutex.Unlock()
		queue, exists := requestTracker[id]
		if exists {
			if len(headers) > 0 {
				queue.headers = headers
			}

			queue.cookies = cookies
			requestTracker[id] = queue
		}
	}

	removeRequestID := func(id proto.NetworkRequestID) {
		requestTrackerMutex.Lock()
		defer requestTrackerMutex.Unlock()
		delete(requestTracker, id)
	}

	var requestWillBeSent proto.NetworkRequestWillBeSent
	_ = p.page.WaitEvent(&requestWillBeSent)

	var requestWillBeSentExtraInfo proto.NetworkRequestWillBeSentExtraInfo
	_ = p.page.WaitEvent(&requestWillBeSentExtraInfo)

	var responseReceived proto.NetworkResponseReceived
	_ = p.page.WaitEvent(&responseReceived)

	var dataReceived proto.NetworkDataReceived
	_ = p.page.WaitEvent(&dataReceived)

	var loadingFailed proto.NetworkLoadingFailed
	_ = p.page.WaitEvent(&loadingFailed)

	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	}

	p.page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent})
	p.page.MustNavigate(url)

	go p.page.EachEvent(func(e *proto.NetworkRequestWillBeSent) {
		logger.Log.Info("Request will be sent", zap.String("url", e.Request.URL), zap.String("type", string(e.Type)))

		headers := make(map[string]interface{})
		for k, v := range e.Request.Headers {
			headers[k] = v.String()
		}

		headersString := utils.ConvertMapToJSON(headers)

		switch e.Type {
		case proto.NetworkResourceTypeDocument:
			redirects = append(redirects, common.Redirect{
				Url:        e.Request.URL,
				StatusCode: 200,
			})

		case proto.NetworkResourceTypeStylesheet:
			addRequestID(e.RequestID, e.Request.URL, utils.TextCSS, headersString)
		case proto.NetworkResourceTypeScript:
			addRequestID(e.RequestID, e.Request.URL, utils.ApplicationJavascript, headersString)
		case proto.NetworkResourceTypeXHR:
			addRequestID(e.RequestID, e.Request.URL, utils.XHR, headersString)
		}
	}, func(e *proto.NetworkResponseReceived) {
		logger.Log.Info("Response received", zap.String("url", e.Response.URL), zap.Int("status", e.Response.Status))

		editStatus(e.RequestID, e.Response.Status)
	}, func(e *proto.NetworkDataReceived) {
		logger.Log.Info("Data received", zap.String("requestId", string(e.RequestID)), zap.Int("length", e.DataLength))

		requestTrackerMutex.Lock()
		queue, exists := requestTracker[e.RequestID]
		requestTrackerMutex.Unlock()

		if exists {
			networkFiles = append(networkFiles, common.FileRequest{
				Src:      queue.src,
				Headers:  queue.headers,
				Cookies:  queue.cookies,
				FileSize: uint(e.DataLength),
				Content:  string(e.Data),
				FileType: string(queue.fileType),
				Duration: queue.duration,
				Status:   queue.status,
			})

			removeRequestID(e.RequestID)
		}
	}, func(e *proto.NetworkRequestWillBeSentExtraInfo) {
		logger.Log.Info("Request will be sent", zap.String("url", string(e.RequestID)))

		headers := make(map[string]interface{})
		for k, v := range e.Headers {
			headers[k] = v.String()
		}

		headersString := utils.ConvertMapToJSON(headers)

		cookies := make(map[string]interface{})
		for k, v := range e.Headers {
			cookies[k] = v.String()
		}

		cookiesString := utils.ConvertMapToJSON(cookies)
		editRequest(e.RequestID, headersString, cookiesString)

	}, func(e *proto.NetworkLoadingFailed) {
		logger.Log.Info("Loading failed", zap.String("type", string(e.Type)), zap.String("text", e.ErrorText))

		removeRequestID(e.RequestID)
	})()

	time.Sleep(3 * time.Second)

	htmlContent, err := p.page.HTML()
	if err != nil {
		return nil, err
	}

	parsedHTML, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	requestTrackerMutex.Lock()
	requestTracker = make(map[proto.NetworkRequestID]TrackerQueue)
	requestTrackerMutex.Unlock()

	return &WebsiteResponse{
		Redirects:  redirects,
		Files:      networkFiles,
		FinalHTML:  htmlContent,
		ParsedHTML: parsedHTML,
		StatusCode: 200,
	}, nil
}

// AnalyzeHTML extracts scripts, styles, and metadata from the parsed HTML.
func (p *PageAnalysis) analyzeHTML(response *WebsiteResponse) (common.WebsiteAnalysis, error) {
	var extractedFiles []common.FileRequest
	title := utils.ExtractTitle(response.ParsedHTML)

	utils.TraverseHTML(response.ParsedHTML, func(node *html.Node) {
		switch node.Data {
		case "script":
			extractedFiles = append(extractedFiles, utils.ProcessScriptNode(node))
		case "style":
			extractedFiles = append(extractedFiles, utils.ProcessStyleNode(node))
		case "link":
			if utils.IsStylesheet(node) {
				extractedFiles = append(extractedFiles, utils.ProcessLinkNode(node))
			} else if utils.IsFont(node) {
				extractedFiles = append(extractedFiles, utils.ProcessFontNode(node))
			}
		}
	})

	extractedFiles = append(extractedFiles, response.Files...)

	return common.WebsiteAnalysis{
		Url:       response.FinalHTML,
		Title:     title,
		Assets:    extractedFiles,
		Redirects: response.Redirects,
	}, nil
}

// logErrorHandling logs errors for better debugging.
func (p *PageAnalysis) logErrorHandling(err error) {
	if strings.Contains(err.Error(), "net::ERR_NAME_NOT_RESOLVED") || strings.Contains(err.Error(), "net::ERR_CONNECTION_REFUSED") {
		logger.Log.Error("Network error", zap.Error(err))
	} else {
		logger.Log.Error("Unknown error", zap.Error(err))
	}
}

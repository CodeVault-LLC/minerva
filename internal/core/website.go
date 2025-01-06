package core

import (
	"encoding/json"
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

	addRequestID := func(id proto.NetworkRequestID, src string, fileType utils.FileType, status int) {
		requestTrackerMutex.Lock()
		defer requestTrackerMutex.Unlock()
		requestTracker[id] = TrackerQueue{
			src:      src,
			fileType: fileType,
			status:   status,
		}
	}

	removeRequestID := func(id proto.NetworkRequestID) {
		requestTrackerMutex.Lock()
		defer requestTrackerMutex.Unlock()
		delete(requestTracker, id)
	}

	var responseReceived proto.NetworkResponseReceived
	_ = p.page.WaitEvent(&responseReceived)

	var requestWillBeSent proto.NetworkRequestWillBeSentExtraInfo
	_ = p.page.WaitEvent(&requestWillBeSent)

	var dataReceived proto.NetworkDataReceived
	_ = p.page.WaitEvent(&dataReceived)

	var loadingFailed proto.NetworkLoadingFailed
	_ = p.page.WaitEvent(&loadingFailed)

	p.page.MustNavigate(url)

	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	}
	p.page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent})

	go p.page.EachEvent(func(e *proto.NetworkResponseReceived) {
		logger.Log.Info("Response received", zap.String("url", e.Response.URL), zap.Int("status", e.Response.Status))

		switch e.Type {
		case proto.NetworkResourceTypeDocument:
			redirects = append(redirects, common.Redirect{
				Url:        e.Response.URL,
				StatusCode: e.Response.Status,
			})
		case proto.NetworkResourceTypeStylesheet:
			addRequestID(e.RequestID, e.Response.URL, utils.TextCSS, e.Response.Status)
		case proto.NetworkResourceTypeScript:
			addRequestID(e.RequestID, e.Response.URL, utils.ApplicationJavascript, e.Response.Status)
		case proto.NetworkResourceTypeXHR:
			addRequestID(e.RequestID, e.Response.URL, utils.XHR, e.Response.Status)
		}
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

		requestTrackerMutex.Lock()
		queue, exists := requestTracker[e.RequestID]
		requestTrackerMutex.Unlock()

		if exists {
			logger.Log.Info("Request Headers", zap.Any("headers", e.Headers), zap.Any("cookies", e.AssociatedCookies))

			headers := make(map[string]string)
			headersJSON, err := json.Marshal(e.Headers)
			if err != nil {
				logger.Log.Error("Failed to marshal headers", zap.Error(err))
				return
			}
			_ = json.Unmarshal(headersJSON, &headers)

			cookies := make(map[string]string)
			cookiesJSON, err := json.Marshal(e.AssociatedCookies)
			if err != nil {
				logger.Log.Error("Failed to marshal cookies", zap.Error(err))
				return
			}
			_ = json.Unmarshal(cookiesJSON, &cookies)

			requestTrackerMutex.Lock()
			queue.headers = string(headersJSON)
			queue.cookies = string(cookiesJSON)
			queue.duration = int(e.ConnectTiming.RequestTime)

			requestTracker[e.RequestID] = queue
			requestTrackerMutex.Unlock()

			removeRequestID(e.RequestID)
		}

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

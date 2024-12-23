package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/codevault-llc/minerva/internal/common"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/responder"
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

func (p *PageAnalysis) FetchWebsite(url, userAgent string) (*WebsiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	go p.page.EachEvent(func(e *proto.NetworkResponseReceived) {
		logger.Log.Info("Response received", zap.String("url", e.Response.URL), zap.String("type", string(e.Type)))
		logger.Log.Info("Response body", zap.String("body", e.Response.))
	})()

	waitForPageLoadEvent := p.page.EachEvent(func(e *proto.PageLoadEventFired) (stop bool) {
		return true
	})
	err := p.page.Navigate(url)
	if err != nil {
		p.logErrorHandling(err)
		return nil, responder.CreateError(responder.ErrInvalidRequest).Error
	}
	waitForPageLoadEvent()

	// Set a realistic User-Agent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
	}

	p.page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent})

	// Initialize request interception
	var redirects []common.Redirect
	var networkFiles []common.FileRequest
	//p.setupRequestInterception(p.page, &redirects, &networkFiles)

	// Wait for page load
	if err := rod.Try(func() {
		p.page.Context(ctx).MustWaitLoad()
	}); err != nil {
		return nil, errors.New("page load timeout or error: " + err.Error())
	}

	// Allow lazy-loaded content to load
	time.Sleep(3 * time.Second)

	htmlContent, err := p.page.HTML()
	if err != nil {
		return nil, err
	}

	parsedHTML, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	return &WebsiteResponse{
		Redirects:  redirects,
		Files:      networkFiles,
		FinalHTML:  htmlContent,
		ParsedHTML: parsedHTML,
		StatusCode: redirects[len(redirects)-1].StatusCode,
	}, nil
}

// setupRequestInterception handles intercepted requests and categorizes responses.
func (p *PageAnalysis) setupRequestInterception(page *rod.Page, redirects *[]common.Redirect, networkFiles *[]common.FileRequest) {
	router := page.HijackRequests()

	router.MustAdd("*", func(c *rod.Hijack) {
		//logger.Log.Info("Request intercepted", zap.String("url", c.Request.URL().String()), zap.String("type", string(c.Request.Type())))
		requestURL := c.Request.URL().String()

		handleRequest := func() {
			switch c.Request.Type() {
			case proto.NetworkResourceTypeDocument:
				if err := rod.Try(func() {
					c.MustLoadResponse()
				}); err != nil {
					logger.Log.Error("Failed to load document response", zap.Error(err), zap.String("url", requestURL))
				}
				*redirects = append(*redirects, common.Redirect{
					Url:        requestURL,
					StatusCode: c.Response.RawResponse.StatusCode,
				})
			case proto.NetworkResourceTypeScript:
				p.processResource(c, string(utils.ApplicationJavascript), requestURL, networkFiles)
			case proto.NetworkResourceTypeStylesheet:
				p.processResource(c, string(utils.TextCSS), requestURL, networkFiles)
			case proto.NetworkResourceTypeFont:
				p.processResource(c, string(utils.Font), requestURL, networkFiles)
			case proto.NetworkResourceTypeXHR:
				p.processResource(c, string(utils.XHR), requestURL, networkFiles)
			}
		}

		handleRequest()
		c.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()
}

// processResource handles intercepted resources, loads their responses, and appends them to the file list.
func (p *PageAnalysis) processResource(c *rod.Hijack, fileType string, requestURL string, networkFiles *[]common.FileRequest) {
	timeStart := time.Now()
	if err := rod.Try(func() {
		c.MustLoadResponse()
	}); err != nil {
		logger.Log.Error("Failed to load resource", zap.Error(err), zap.String("url", requestURL))
	}
	*networkFiles = append(*networkFiles, common.FileRequest{
		Src:        requestURL,
		Content:    c.Response.Body(),
		HashedBody: utils.SHA256(c.Response.Body()),
		FileSize:   uint(len(c.Response.Body())),
		FileType:   fileType,
		Duration:   int(time.Since(timeStart).Milliseconds()),
	})
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

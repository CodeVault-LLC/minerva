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

var browser *rod.Browser

// InitializeBrowser sets up the rod browser instance.
func InitializeBrowser() {
	if browser == nil {
		browser = rod.New().MustConnect().NoDefaultDevice().MustIncognito()
	}
}

// CloseBrowser terminates the rod browser instance.
func CloseBrowser() {
	if browser != nil {
		browser.MustClose()
	}
}

type WebsiteResponse struct {
	Redirects    []common.Redirect
	Files        []common.FileRequest
	FinalHTML    string
	ParsedHTML   *html.Node
	WebsiteTitle string
}

// FetchWebsite retrieves the website content and its network resources.
func FetchWebsite(url, userAgent string) (*WebsiteResponse, error) {
	ctx := context.Background()

	page, err := browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		return nil, responder.CreateError(responder.ErrInvalidRequest).Error
	}
	defer page.Close()

	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: userAgent,
	})

	var redirects []common.Redirect
	var networkFiles []common.FileRequest

	router := page.HijackRequests()

	router.MustAdd("*.*", func(c *rod.Hijack) {
		requestURL := c.Request.URL().String()

		handleRequest := func() {
			switch c.Request.Type() {
			case proto.NetworkResourceTypeScript:
				startTime := time.Now()

				if err := rod.Try(func() {
					c.MustLoadResponse()
				}); err != nil {
					logger.Log.Error("Failed to load script response", zap.Error(err), zap.String("url", requestURL))
				}

				networkFiles = append(networkFiles, common.FileRequest{
					Src:        requestURL,
					Content:    c.Response.Body(),
					HashedBody: utils.SHA256(c.Response.Body()),
					FileSize:   uint(len(c.Response.Body())),
					FileType:   string(utils.ApplicationJavascript),
					Duration:   int(time.Since(startTime).Milliseconds()),
				})
			/*case proto.NetworkResourceTypeDocument:
			if err := rod.Try(func() {
				c.MustLoadResponse()
			}); err != nil {
				logger.Log.Error("Failed to load css response", zap.Error(err), zap.String("url", requestURL))
			}

			logger.Log.Info("Document request intercepted", zap.String("url", requestURL))

			redirects = append(redirects, common.Redirect{
				Url: requestURL,
				Screenshot: common.Screenshot{
					Content: string(page.MustWaitStable().MustScreenshotFullPage()),
				},
				StatusCode: c.Response.RawResponse.StatusCode,
			})*/
			case proto.NetworkResourceTypeStylesheet:
				timeStart := time.Now()

				if err := rod.Try(func() {
					c.MustLoadResponse()
				}); err != nil {
					logger.Log.Error("Failed to load css response", zap.Error(err), zap.String("url", requestURL))
				}

				networkFiles = append(networkFiles, common.FileRequest{
					Src:        requestURL,
					Content:    c.Response.Body(),
					HashedBody: utils.SHA256(c.Response.Body()),
					FileSize:   uint(len(c.Response.Body())),
					FileType:   string(utils.TextCSS),
					Duration:   int(time.Since(timeStart).Milliseconds()),
				})
			case proto.NetworkResourceTypeFont:
				timeStart := time.Now()

				if err := rod.Try(func() {
					c.MustLoadResponse()
				}); err != nil {
					logger.Log.Error("Failed to load font response", zap.Error(err), zap.String("url", requestURL))
				}

				networkFiles = append(networkFiles, common.FileRequest{
					Src:        requestURL,
					Content:    c.Response.Body(),
					HashedBody: utils.SHA256(c.Response.Body()),
					FileSize:   uint(len(c.Response.Body())),
					FileType:   string(utils.Font),
					Duration:   int(time.Since(timeStart).Milliseconds()),
				})
				/*case proto.NetworkResourceTypeXHR:
					if err := rod.Try(func() {
						c.MustLoadResponse()
					}); err != nil {
						logger.Log.Error("Failed to load xhr response", zap.Error(err), zap.String("url", requestURL))
					}

					networkFiles = append(networkFiles, common.FileRequest{
						Src:        requestURL,
						Content:    c.Response.Body(),
						HashedBody: utils.SHA256(c.Response.Body()),
						FileSize:   uint(len(c.Response.Body())),
						FileType:   string(utils.XHR),
					})
				case proto.NetworkResourceTypeWebSocket:
					if err := rod.Try(func() {
						c.MustLoadResponse()
					}); err != nil {
						logger.Log.Error("Failed to load websocket response", zap.Error(err), zap.String("url", requestURL))
					}
				*/
			}
		}

		handleRequest()

		c.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	if err := rod.Try(func() {
		page.Context(ctx).MustWaitLoad()
	}); err != nil {
		return nil, errors.New("page load timeout or error: " + err.Error())
	}

	htmlContent, err := page.HTML()
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
	}, nil
}

// AnalyzeHTML extracts scripts, styles, and metadata from the parsed HTML.
func AnalyzeHTML(response *WebsiteResponse) (common.WebsiteAnalysis, error) {
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
		Url:        response.FinalHTML,
		Title:      title,
		StatusCode: 200,
		Assets:     extractedFiles,
		Redirects:  response.Redirects,
	}, nil
}

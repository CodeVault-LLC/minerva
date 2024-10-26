package websites

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
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
	Redirects    []string
	Files        []types.FileRequest
	FinalHTML    string
	ParsedHTML   *html.Node
	WebsiteTitle string
}

// FetchWebsite retrieves the website content and its network resources.
func FetchWebsite(url, userAgent string) (*WebsiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	page := browser.MustPage(url)
	defer page.Close()

	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: userAgent,
	})

	var redirects []string
	var networkFiles []types.FileRequest

	router := page.HijackRequests()

	router.MustAdd("*.*", func(c *rod.Hijack) {
		requestURL := c.Request.URL().String()

		logger.Log.Info("Request intercepted", zap.String("url", requestURL))

		handleRequest := func() {
			switch c.Request.Type() {
			case proto.NetworkResourceTypeScript:
				if err := rod.Try(func() {
					c.MustLoadResponse()
				}); err != nil {
					logger.Log.Error("Failed to load script response", zap.Error(err), zap.String("url", requestURL))
				}

				networkFiles = append(networkFiles, types.FileRequest{
					Src:        requestURL,
					Content:    c.Response.Body(),
					HashedBody: utils.SHA256(c.Response.Body()),
					FileSize:   uint(len(c.Response.Body())),
					FileType:   "application/javascript",
				})
			case proto.NetworkResourceTypeDocument:
				redirects = append(redirects, requestURL)
			case proto.NetworkResourceTypeStylesheet:
				if err := rod.Try(func() {
					c.MustLoadResponse()
				}); err != nil {
					logger.Log.Error("Failed to load css response", zap.Error(err), zap.String("url", requestURL))
				}

				networkFiles = append(networkFiles, types.FileRequest{
					Src:        requestURL,
					Content:    c.Response.Body(),
					HashedBody: utils.SHA256(c.Response.Body()),
					FileSize:   uint(len(c.Response.Body())),
					FileType:   "text/css",
				})
			case proto.NetworkResourceTypeFont:
				if err := rod.Try(func() {
					c.MustLoadResponse()
				}); err != nil {
					logger.Log.Error("Failed to load font response", zap.Error(err), zap.String("url", requestURL))
				}

				networkFiles = append(networkFiles, types.FileRequest{
					Src:        requestURL,
					Content:    c.Response.Body(),
					HashedBody: utils.SHA256(c.Response.Body()),
					FileSize:   uint(len(c.Response.Body())),
					FileType:   "font",
				})
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
func AnalyzeHTML(response *WebsiteResponse) (types.WebsiteAnalysis, error) {
	var extractedFiles []types.FileRequest
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

	return types.WebsiteAnalysis{
		Url:        response.FinalHTML,
		Title:      title,
		StatusCode: 200,
		Assets:     extractedFiles,
		Redirects:  response.Redirects,
	}, nil
}

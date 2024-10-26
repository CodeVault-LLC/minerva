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
	title := extractTitle(response.ParsedHTML)

	traverseHTML(response.ParsedHTML, func(node *html.Node) {
		switch node.Data {
		case "script":
			extractedFiles = append(extractedFiles, processScriptNode(node))
		case "style":
			extractedFiles = append(extractedFiles, processStyleNode(node))
		case "link":
			if isStylesheet(node) {
				extractedFiles = append(extractedFiles, processLinkNode(node))
			} else if isFont(node) {
				extractedFiles = append(extractedFiles, processFontNode(node))
			}
		}
	})

	// Combine extracted inline scripts/styles and those gathered during the network requests.
	extractedFiles = append(extractedFiles, response.Files...)

	return types.WebsiteAnalysis{
		Url:        response.FinalHTML,
		Title:      title,
		StatusCode: 200,
		Assets:     extractedFiles,
		Redirects:  response.Redirects,
	}, nil
}

// extractTitle retrieves the title from the parsed HTML.
func extractTitle(doc *html.Node) string {
	var title string
	traverseHTML(doc, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" && node.FirstChild != nil {
			title = node.FirstChild.Data
		}
	})
	return title
}

// processScriptNode extracts data from a script element.
func processScriptNode(node *html.Node) types.FileRequest {
	// Handle inline script.
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		return createFileRequest("inline-script", node.FirstChild.Data, "application/javascript")
	}
	return types.FileRequest{}
}

// processStyleNode extracts data from a style element.
func processStyleNode(node *html.Node) types.FileRequest {
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		return createFileRequest("inline-style", node.FirstChild.Data, "text/css")
	}
	return types.FileRequest{}
}

// processLinkNode extracts data from a link element if it's a stylesheet.
func processLinkNode(node *html.Node) types.FileRequest {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return types.FileRequest{
				Src:      attr.Val,
				FileType: "text/css",
			}
		}
	}
	return types.FileRequest{}
}

// processFontNode extracts data from a link element if it's a font.
func processFontNode(node *html.Node) types.FileRequest {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return types.FileRequest{
				Src:      attr.Val,
				FileType: "font",
			}
		}
	}
	return types.FileRequest{}
}

// isStylesheet checks if a link element is a stylesheet.
func isStylesheet(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "rel" && attr.Val == "stylesheet" {
			return true
		}
	}
	return false
}

// isFont checks if a link element points to a font.
func isFont(node *html.Node) bool {
	for _, attr := range node.Attr {
		if attr.Key == "rel" && strings.Contains(attr.Val, "font") {
			return true
		}
	}
	return false
}

// createFileRequest constructs a FileRequest with content details.
func createFileRequest(src, content, fileType string) types.FileRequest {
	return types.FileRequest{
		Src:        src,
		Content:    content,
		HashedBody: utils.SHA256(content),
		FileSize:   uint(len(content)),
		FileType:   fileType,
	}
}

// traverseHTML recursively walks through the HTML nodes and applies the given function.
func traverseHTML(node *html.Node, fn func(*html.Node)) {
	fn(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverseHTML(child, fn)
	}
}

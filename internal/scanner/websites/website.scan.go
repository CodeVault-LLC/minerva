package websites

import (
	"context"
	"errors"
	"io"
	"net/http"
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
		browser = rod.New().MustConnect().NoDefaultDevice()
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
	Scripts      []types.FileRequest
	FinalHTML    string
	ParsedHTML   *html.Node
	WebsiteTitle string
}

// FetchWebsite retrieves the website content and its scripts.
func FetchWebsite(url, userAgent string) (*WebsiteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page := browser.MustPage(url)
	defer page.Close()

	// Set the user agent before navigation
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: userAgent,
	})

	var redirects []string
	var scriptFiles []types.FileRequest

	// Capture network events to track redirects and script URLs.
	go page.EachEvent(func(event *proto.NetworkRequestWillBeSent) {
		if event.Type == proto.NetworkResourceTypeDocument && event.Request.URL != url {
			redirects = append(redirects, event.Request.URL)
		} else if event.Type == proto.NetworkResourceTypeScript && event.Request.URL != "" {
			content, err := downloadContent(event.Request.URL)
			if err != nil {
				logger.Log.Error("Failed to download script", zap.Error(err))
			}
			scriptFiles = append(scriptFiles, createFileRequest(event.Request.URL, content, "application/javascript"))
		}
	})()

	// Navigate and wait for the page to load.
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
		Scripts:    scriptFiles,
		FinalHTML:  htmlContent,
		ParsedHTML: parsedHTML,
	}, nil
}

// AnalyzeHTML extracts scripts, styles, and metadata from the parsed HTML.
func AnalyzeHTML(response *WebsiteResponse) (types.WebsiteAnalysis, error) {
	var scriptsAndStyles []types.FileRequest
	title := extractTitle(response.ParsedHTML)

	traverseHTML(response.ParsedHTML, func(node *html.Node) {
		switch node.Data {
		case "script":
			scriptsAndStyles = append(scriptsAndStyles, processScriptNode(node))
		case "style":
			scriptsAndStyles = append(scriptsAndStyles, processStyleNode(node))
		case "link":
			if isStylesheet(node) {
				scriptsAndStyles = append(scriptsAndStyles, processLinkNode(node))
			}
		}
	})

	// Combine extracted scripts and those gathered during the network requests.
	scriptsAndStyles = append(scriptsAndStyles, response.Scripts...)

	return types.WebsiteAnalysis{
		Url:        response.FinalHTML,
		Title:      title,
		StatusCode: 200,
		Assets:     scriptsAndStyles,
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
	for _, attr := range node.Attr {
		if attr.Key == "src" {
			url := attr.Val
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}
			content, err := downloadContent(url)
			if err != nil {
				logger.Log.Error("Failed to download external script", zap.String("src", attr.Val), zap.Error(err))
				return types.FileRequest{}
			}
			return createFileRequest(attr.Val, content, "application/javascript")
		}
	}

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
			url := attr.Val
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}

			content, err := downloadContent(url)
			if err != nil {
				logger.Log.Error("Failed to download stylesheet", zap.String("href", attr.Val), zap.Error(err))
				return types.FileRequest{}
			}
			return createFileRequest(attr.Val, content, "text/css")
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

// downloadContent retrieves content from a URL.
func downloadContent(url string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to download content, status code: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// traverseHTML recursively walks through the HTML nodes and applies the given function.
func traverseHTML(node *html.Node, fn func(*html.Node)) {
	fn(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverseHTML(child, fn)
	}
}

package websites

import (
	"context"
	"strings"
	"time"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"golang.org/x/net/html"
)

var browser *rod.Browser

// Initialize the browser instance globally and keep it running
func InitBrowser() {
	if browser == nil {
		browser = rod.New().MustConnect().NoDefaultDevice()
	}
}

// Close the browser instance when the application is shutting down
func CloseBrowser() {
	if browser != nil {
		browser.MustClose()
	}
}

type RequestWebsiteResponse struct {
	RedirectChain   []string
	JavaScriptFiles []string
	FinalHTML       string
	ParsedBody      *html.Node
}

func RequestWebsite(url string, userAgent string) (*RequestWebsiteResponse, error) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page := browser.MustPage(url)
	defer page.Close()

	var redirectChain []string
	var jsFiles []string

	go page.EachEvent(func(e *proto.NetworkRequestWillBeSent) {
		if e.Type == "Document" && e.Request.URL != url {
			redirectChain = append(redirectChain, e.Request.URL)
		}
		if e.Type == "Script" {
			jsFiles = append(jsFiles, e.Request.URL)
		}
	})()

	// Set the user-agent and navigate to the URL
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      userAgent,
		AcceptLanguage: "en-US,en;q=0.9",
		Platform:       "Win32",
	})

	// Wait for the page to load or for the context to time out
	err := rod.Try(func() {
		page.Context(ctx).MustWaitLoad()
	})

	// If the context times out, proceed with whatever data is available
	if err != nil {
		// Log or handle the timeout error if needed
	}

	// Capture the final HTML of the page, ignoring any further wait
	htmlContent, err := page.HTML()
	if err != nil {
		return nil, err
	}

	// Parse the final HTML into a DOM structure
	respBody := strings.NewReader(htmlContent)
	doc, err := html.Parse(respBody)
	if err != nil {
		return nil, err
	}

	return &RequestWebsiteResponse{
		RedirectChain:   redirectChain,
		JavaScriptFiles: jsFiles,
		FinalHTML:       htmlContent,
		ParsedBody:      doc,
	}, nil
}

// AnalyzeWebsite traverses the parsed HTML and extracts information.
func AnalyzeWebsite(resp *RequestWebsiteResponse) (models.ScanResponse, error) {
	var websiteName string = "Unknown"

	// Traverse and extract data from the parsed HTML
	var f func(*html.Node)
	f = func(n *html.Node) {
		// Extract website title if it exists
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			websiteName = n.FirstChild.Data
		}
		// Continue traversing
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// Start traversing from the root of the parsed document
	f(resp.ParsedBody)

	// Return a more flexible response with various details
	return models.ScanResponse{
		WebsiteUrl:  resp.FinalHTML,
		WebsiteName: websiteName,
		StatusCode:  200,
		Javascript:  resp.JavaScriptFiles,
	}, nil
}

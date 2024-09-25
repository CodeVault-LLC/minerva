package websites

import (
	"io"
	"net/http"
	"strings"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"golang.org/x/net/html"
)

type RequestWebsiteResponse struct {
	Response      *http.Response
	ParsedBody    *html.Node
	RedirectChain []string
}

func RequestWebsite(url string, userAgent string) (*RequestWebsiteResponse, error) {
	var redirectChain []string

	// Custom HTTP client with a redirect policy
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Append each redirected URL to the chain
			redirectChain = append(redirectChain, req.URL.String())
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &RequestWebsiteResponse{}, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "en-US")

	resp, err := client.Do(req)
	if err != nil {
		return &RequestWebsiteResponse{}, err
	}

	// Add the final URL in the redirect chain
	redirectChain = append(redirectChain, resp.Request.URL.String())

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RequestWebsiteResponse{}, err
	}

	resp.Body = io.NopCloser(strings.NewReader(string(responseData)))

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return &RequestWebsiteResponse{}, err
	}

	return &RequestWebsiteResponse{
		Response:      resp,
		ParsedBody:    doc,
		RedirectChain: redirectChain, // Save the redirect chain
	}, nil
}

func AnalyzeWebsite(resp *RequestWebsiteResponse) (models.ScanResponse, error) {
	var websiteName string = "Unknown"

	// Function to traverse and extract data
	var f func(*html.Node)
	f = func(n *html.Node) {
		// Extract website name
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			websiteName = n.FirstChild.Data
		}

		// Recursively traverse children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// Start traversing from the root document
	f(resp.ParsedBody)

	return models.ScanResponse{
		WebsiteName: websiteName,
		WebsiteUrl:  resp.Response.Request.URL.Hostname(),
		StatusCode:  resp.Response.StatusCode,
	}, nil
}

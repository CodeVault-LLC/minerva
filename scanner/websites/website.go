package websites

import (
	"io"
	"net/http"
	"strings"

	"github.com/codevault-llc/humblebrag-api/models"
	"golang.org/x/net/html"
)

func ScanWebsite(url string) (models.ScanResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.ScanResponse{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return models.ScanResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return models.ScanResponse{}, err
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ScanResponse{}, err
	}

	resp.Body = io.NopCloser(strings.NewReader(string(responseData)))

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return models.ScanResponse{}, err
	}

	var websiteName string = "Unknown"
	var extractedScripts []models.ScriptRequest

	// Function to traverse and extract data
	var f func(*html.Node)
	f = func(n *html.Node) {
		// Extract website name
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			websiteName = n.FirstChild.Data
		}

		// Extract script tags
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					// Handle external scripts
					scriptContent, err := downloadExternalScript(a.Val)
					if err != nil {
						return
					}
					extractedScripts = append(extractedScripts, models.ScriptRequest{
						Src:     a.Val,
						Content: scriptContent,
					})
				}
			}

			// Handle inline scripts
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				extractedScripts = append(extractedScripts, models.ScriptRequest{
					Src:     "inline-script",
					Content: n.FirstChild.Data,
				})
			}
		}

		// Recursively traverse children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	// Start traversing from the root document
	f(doc)

	return models.ScanResponse{
		WebsiteName: websiteName,
		Scripts:     extractedScripts,
		WebsiteUrl:  url,
	}, nil
}

func downloadExternalScript(url string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set config
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", err
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

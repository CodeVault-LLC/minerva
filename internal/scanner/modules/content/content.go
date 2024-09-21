package content

import (
	"io"
	"net/http"
	"sync"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
	"golang.org/x/net/html"
)

func ContentModule(scanId uint, requestedWebsite *html.Node) {
	scripts, err := findScripts(requestedWebsite)
	if err != nil {
		logger.Log.Error("Failed to find scripts: %v", err)
		return
	}

	for _, script := range scripts {
		content := models.ContentModel{
			ScanID: scanId,

			Name:    script.Src,
			Content: script.Content,
		}

		_, err := service.CreateContent(content)
		if err != nil {
			logger.Log.Error("Failed to save content: %v", err)
		}
	}

	findings := scanSecrets(scripts)

	service.CreateFindings(scanId, findings)
}

func findScripts(doc *html.Node) ([]models.ScriptRequest, error) {
	var extractedScripts []models.ScriptRequest

	var f func(*html.Node)
	f = func(n *html.Node) {

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

	return extractedScripts, nil
}

func downloadExternalScript(url string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set config
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")

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

func scanSecrets(scripts []models.ScriptRequest) []utils.RegexReturn {
	var results []utils.RegexReturn

	var wg sync.WaitGroup
	var mu sync.Mutex

	concurrencyLimit := make(chan struct{}, 10)

	for _, rule := range config.ConfigRules {
		concurrencyLimit <- struct{}{}
		wg.Add(1)

		go func(rule types.Rule) {
			defer wg.Done()
			defer func() { <-concurrencyLimit }()

			var scriptResults []utils.Match
			for _, script := range scripts {
				matches := utils.GenericScan(rule, utils.Script(script))
				if len(matches) > 0 {
					scriptResults = append(scriptResults, matches...)
				}
			}

			if len(scriptResults) > 0 {
				mu.Lock()
				results = append(results, utils.RegexReturn{Name: rule.RuleID, Matches: scriptResults, Description: rule.Description})
				mu.Unlock()
			}
		}(*rule)
	}

	wg.Wait()
	return results
}

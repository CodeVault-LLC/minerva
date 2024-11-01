package modules

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/repository"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"golang.org/x/net/html"
)

type MetadataModule struct{}

func NewMetadataModule() *MetadataModule {
	return &MetadataModule{}
}

func (m *MetadataModule) Execute(job entities.JobModel, website types.WebsiteAnalysis) error {
	client := &http.Client{Timeout: 15 * time.Second}

	robotsTxtChan := make(chan string)
	readmeChan := make(chan string)
	licenseChan := make(chan string)
	cmsChan := make(chan string)
	serverSoftwareChan := make(chan string)
	frameworksChan := make(chan []string)
	languageChan := make(chan string)
	headersChan := make(chan map[string]string)

	var wg sync.WaitGroup
	wg.Add(7)

	go func() {
		defer wg.Done()
		robotsTxtChan <- fetchFileContent(client, job.URL+"/robots.txt")
	}()

	go func() {
		defer wg.Done()
		readmeChan <- fetchFileContent(client, job.URL+"/readme.html")
	}()

	go func() {
		defer wg.Done()
		licenseChan <- fetchFileContent(client, job.URL+"/license.txt")
	}()

	go func() {
		defer wg.Done()
		//cmsChan <- detectCMS(client, url)
		cmsChan <- "Unknown"
	}()

	go func() {
		defer wg.Done()
		serverSoftwareChan <- getServerHeader(client, job.URL)
	}()

	go func() {
		defer wg.Done()
		frameworksChan <- detectFrameworks(client, job.URL)
	}()

	go func() {
		defer wg.Done()
		languageChan <- detectLanguage(client, job.URL)
	}()

	go func() {
		wg.Wait()
		close(robotsTxtChan)
		close(readmeChan)
		close(licenseChan)
		close(cmsChan)
		close(serverSoftwareChan)
		close(frameworksChan)
		close(languageChan)
		close(headersChan)
	}()
	resultModel := entities.MetadataModel{
		ScanID:         job.ScanID,
		Robots:         <-robotsTxtChan,
		Readme:         <-readmeChan,
		License:        <-licenseChan,
		CMS:            <-cmsChan,
		ServerSoftware: <-serverSoftwareChan,
		Frameworks:     <-frameworksChan,
		ServerLanguage: <-languageChan,
	}

	_, err := repository.MetadataRepository.Create(resultModel)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetadataModule) Name() string {
	return "Metadata"
}

// fetchFileContent retrieves the content of a file from the given URL.
func fetchFileContent(client *http.Client, url string) string {
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return ""
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(content)
}

// detectCMS improves CMS detection by checking meta tags, common CMS directories, and the /feed endpoint for WordPress.
/*func detectCMS(client *http.Client, url string) string {
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "Unknown"
	}
	defer resp.Body.Close()

	// Tokenize HTML to detect CMS via meta tags

	// Check for CMS-specific URLs or files
	if checkFileExists(client, url+"/wp-content/") || checkFileExists(client, url+"/wp-includes/") {
		//checkWordPressVersion(client, url)
		return "WordPress"
	}
	if checkFileExists(client, url+"/administrator/") {
		return "Joomla"
	}
	if checkFileExists(client, url+"/core/") {
		return "Drupal"
	}

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return "Unknown"
		case html.StartTagToken:
			t := tokenizer.Token()
			if t.Data == "meta" {
				for _, attr := range t.Attr {
					if attr.Key == "name" && (strings.Contains(attr.Val, "generator") || strings.Contains(attr.Val, "cms")) {
						return getMetaContent(t.Attr)
					}
				}
			}
		}
	}
}*/

// checkWordPressVersion checks the WordPress version from the /feed page.
/*func checkWordPressVersion(client *http.Client, url string) string {
	resp, err := client.Get(url + "/feed")
	if err != nil || resp.StatusCode != http.StatusOK {
		return "WordPress (version unknown)"
	}
	defer resp.Body.Close()

	// Check if feed contains version info
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if strings.Contains(bodyString, "wordpress.org/?v=") {
		versionIndex := strings.Index(bodyString, "wordpress.org/?v=")
		version := bodyString[versionIndex+17 : versionIndex+22] // assuming version follows right after
		return fmt.Sprintf("WordPress %s", version)
	}

	return "WordPress (version unknown)"
}*/

// detectLanguage inspects headers or URLs to determine the server-side language (PHP, Python, etc.).
func detectLanguage(client *http.Client, url string) string {
	resp, err := client.Head(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	// Check headers for language hints
	if strings.Contains(resp.Header.Get("X-Powered-By"), "PHP") {
		return "PHP"
	}
	if strings.Contains(resp.Header.Get("X-Powered-By"), "Python") {
		return "Python"
	}
	if strings.Contains(resp.Header.Get("X-Powered-By"), "Ruby") {
		return "Ruby"
	}
	if strings.Contains(resp.Header.Get("X-Powered-By"), "Node.js") {
		return "Node.js"
	}

	// Fallback by checking common file extensions
	if checkFileExists(client, url+"/index.php") {
		return "PHP"
	}
	if checkFileExists(client, url+"/app.py") {
		return "Python"
	}

	return "Unknown"
}

// getServerHeader retrieves the "Server" header to infer the web server software.
func getServerHeader(client *http.Client, url string) string {
	resp, err := client.Head(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	server := resp.Header.Get("Server")
	if server == "" {
		return "Unknown"
	}
	return server
}

// detectFrameworks tries to detect any common JavaScript frameworks or libraries by inspecting the HTML content or headers.
func detectFrameworks(client *http.Client, url string) []string {
	var frameworks []string

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return frameworks
	}
	defer resp.Body.Close()

	// Detect common JavaScript frameworks from headers or content
	if strings.Contains(resp.Header.Get("X-Powered-By"), "Express") {
		frameworks = append(frameworks, "Express.js")
	}

	// Tokenize HTML to look for specific framework-related tags
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		t := tokenizer.Token()
		if t.Data == "script" {
			for _, attr := range t.Attr {
				if attr.Key == "src" && strings.Contains(attr.Val, "jquery") {
					frameworks = append(frameworks, "jQuery")
				}
				if attr.Key == "src" && strings.Contains(attr.Val, "angular") {
					frameworks = append(frameworks, "Angular")
				}
				if attr.Key == "src" && strings.Contains(attr.Val, "react") {
					frameworks = append(frameworks, "React")
				}
			}
		}
	}

	return frameworks
}

// getMetaContent extracts the content attribute from a meta tag.
/*func getMetaContent(attrs []html.Attribute) string {
	for _, attr := range attrs {
		if attr.Key == "content" {
			return attr.Val
		}
	}
	return "Unknown"
}*/

// checkFileExists checks if a file exists on the given URL.
func checkFileExists(client *http.Client, url string) bool {
	resp, err := client.Head(url)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}

package scanner

import (
	"fmt"
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/content"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/list"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/metadata"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/network"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/nmap"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"golang.org/x/net/html"
)

func ScanWebsite(url string, licenseId uint) (models.ScanModel, error) {
	logger.Log.Info("Starting website scan for URL: %s", url)

	// Initial website scan
	requestedWebsite, err := websites.RequestWebsite(url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	if err != nil {
		logger.Log.Error("Failed to scan website: %v", err)
		return models.ScanModel{}, err
	}

	website, err := websites.AnalyzeWebsite(requestedWebsite)
	if err != nil {
		logger.Log.Error("Failed to analyze website: %v", err)
		return models.ScanModel{}, err
	}

	// Save initial scan result
	scanModel := models.ScanModel{
		Url:           url,
		Name:          website.WebsiteName,
		RedirectChain: requestedWebsite.RedirectChain,
		StatusCode:    website.StatusCode,
		Status:        models.ScanStatusPending, // Set status to pending for further processing
		LicenseID:     licenseId,
		Sha256:        utils.SHA256(website.WebsiteUrl),
		SHA1:          utils.SHA1(website.WebsiteUrl),
		MD5:           utils.MD5(website.WebsiteUrl),
	}

	scanModel, err = service.CreateScan(scanModel)
	if err != nil {
		logger.Log.Error("Failed to save initial scan result: %v", err)
		return models.ScanModel{}, err
	}

	// Start background goroutines to handle further scans asynchronously
	go runBackgroundModules(scanModel.ID, url, requestedWebsite.ParsedBody)

	// Return the response immediately, without waiting for the background tasks
	return scanModel, nil
}

func runBackgroundModules(scanId uint, url string, requestedWebsite *html.Node) {
	var wg sync.WaitGroup
	wg.Add(5) // Amount of modules

	go func() {
		fmt.Println("Starting content module")
		defer wg.Done() // Signal that the content module is done
		content.ContentModule(scanId, requestedWebsite)
	}()

	go func() {
		fmt.Println("Starting list module")
		defer wg.Done() // Signal that the list module is done
		list.ListModule(scanId, url)
	}()

	go func() {
		fmt.Println("Starting metadata module")
		defer wg.Done() // Signal that the list module is done
		metadata.MetadataModule(scanId, url)
	}()

	go func() {
		fmt.Println("Starting nmap module")
		defer wg.Done()
		nmap.NmapModule(scanId, url)
	}()

	go func() {
		fmt.Println("Starting network module")
		defer wg.Done() // Signal that the network module is done
		network.NetworkModule(scanId, url)
	}()

	// Wait for all background processes to complete
	go func() {
		wg.Wait() // This will block until all 3 modules call `wg.Done()`
		fmt.Println("All modules are done")

		// Once all modules are done, update the scan status to "Complete"
		scanModel := models.ScanModel{
			Status: models.ScanStatusComplete,
		}
		scanModel.ID = scanId

		_, err := service.UpdateScan(scanModel)
		if err != nil {
			logger.Log.Error("Failed to update scan status: %v", err)
		}
	}()
}

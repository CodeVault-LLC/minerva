package scanner

import (
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/content"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/list"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/network"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"golang.org/x/net/html"
)

func ScanWebsite(url string, userId uint) (models.ScanModel, error) {
	logger.Log.Info("Starting website scan for URL: %s", url)

	// Initial website scan
	requestedWebsite, err := websites.RequestWebsite(url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	if err != nil {
		logger.Log.Error("Failed to scan website: %v", err)
		return models.ScanModel{}, err
	}

	logger.Log.Info("Redirect chain: %v", requestedWebsite.RedirectChain)

	website, err := websites.AnalyzeWebsite(requestedWebsite)
	if err != nil {
		logger.Log.Error("Failed to analyze website: %v", err)
		return models.ScanModel{}, err
	}

	// Save initial scan result
	scanModel := models.ScanModel{
		WebsiteUrl:    url,
		WebsiteName:   website.WebsiteName,
		RedirectChain: requestedWebsite.RedirectChain,
		StatusCode:    website.StatusCode,
		Status:        models.ScanStatusPending, // Set status to pending for further processing
		UserID:        userId,
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
	logger.Log.Info("Running background modules for scan ID: %d", scanId)

	var wg sync.WaitGroup
	wg.Add(3) // Adding 3 because we have 3 modules

	go func() {
		defer wg.Done() // Signal that the network module is done
		logger.Log.Info("Starting network module for scan ID: %d", scanId)
		network.NetworkModule(scanId, url)
		logger.Log.Info("Completed network module for scan ID: %d", scanId)
	}()

	go func() {
		defer wg.Done() // Signal that the content module is done
		logger.Log.Info("Starting content module for scan ID %d", scanId)
		content.ContentModule(scanId, requestedWebsite)
		logger.Log.Info("Completed content module for scan ID %d", scanId)
	}()

	go func() {
		defer wg.Done() // Signal that the list module is done
		logger.Log.Info("Starting list module for scan ID %d", scanId)
		list.ListModule(scanId, url)
		logger.Log.Info("Completed list module for scan ID %d", scanId)
	}()

	// Wait for all background processes to complete
	go func() {
		wg.Wait() // This will block until all 3 modules call `wg.Done()`

		// Once all modules are done, update the scan status to "Complete"
		logger.Log.Info("All background modules completed for scan ID: %d", scanId)
		scanModel := models.ScanModel{
			Status: models.ScanStatusComplete,
		}
		scanModel.ID = scanId

		_, err := service.UpdateScan(scanModel)
		if err != nil {
			logger.Log.Error("Failed to update scan status: %v", err)
		} else {
			logger.Log.Info("Successfully updated scan status to 'Complete' for scan ID: %d", scanId)
		}
	}()
}

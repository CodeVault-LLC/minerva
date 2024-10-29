package scanner

import (
	"fmt"
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/content"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/list"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/metadata"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/network"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"go.uber.org/zap"
)

func ScanWebsite(url string, userAgent string, licenseId uint) (models.ScanModel, error) {
	logger.Log.Info("Starting website scan for URL: %s", zap.String("url", url))

	// Initial website scan
	requestedWebsite, err := websites.FetchWebsite(url, userAgent)
	if err != nil {
		return models.ScanModel{}, err
	}

	website, err := websites.AnalyzeHTML(requestedWebsite)
	if err != nil {
		logger.Log.Error("Failed to analyze website: %v", zap.Error(err))
		return models.ScanModel{}, err
	}

	// Save initial scan result
	scanModel := models.ScanModel{
		Url:           url,
		Title:         website.Title,
		RedirectChain: requestedWebsite.Redirects,
		StatusCode:    website.StatusCode,
		Status:        models.ScanStatusPending, // Set status to pending for further processing
		LicenseID:     licenseId,
		Sha256:        utils.SHA256(website.Url),
		SHA1:          utils.SHA1(website.Url),
		MD5:           utils.MD5(website.Url),
	}

	scanModel, err = service.CreateScan(scanModel)
	if err != nil {
		logger.Log.Error("Failed to save initial scan result: %v", zap.Error(err))
		return models.ScanModel{}, err
	}

	// Start background goroutines to handle further scans asynchronously
	go runBackgroundModules(scanModel.ID, url, website)

	// Return the response immediately, without waiting for the background tasks
	return scanModel, nil
}

func runBackgroundModules(scanId uint, url string, website types.WebsiteAnalysis) {
	var wg sync.WaitGroup
	wg.Add(4) // Amount of modules

	go func() {
		fmt.Println("Starting content module")
		defer wg.Done() // Signal that the content module is done
		content.ContentModule(scanId, website.Assets)
	}()

	go func() {
		fmt.Println("Starting list module")
		defer wg.Done() // Signal that the list module is done
		list.ListModule(scanId, url)
	}()

	go func() {
		fmt.Println("Starting metadata module")
		defer wg.Done() // Signal that the list module is done
		_, err := metadata.MetadataModule(scanId, url)
		if err != nil {
			logger.Log.Error("Failed to run metadata module: %v", zap.Error(err))
		}
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
			logger.Log.Error("Failed to update scan status: %v", zap.Error(err))
		}
	}()
}

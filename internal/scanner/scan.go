package scanner

import (
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/content"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/finding"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/list"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/modules/network"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func ScanWebsite(url string, userId uint) (models.ScanModel, error) {
	logger.Log.Info("Starting website scan for URL: %s", url)

	// Initial website scan
	website, err := websites.ScanWebsite(url)
	if err != nil {
		logger.Log.Error("Failed to scan website: %v", err)
		return models.ScanModel{}, err
	}

	// Save initial scan result
	scanModel := models.ScanModel{
		WebsiteUrl:  website.WebsiteUrl,
		WebsiteName: website.WebsiteName,
		Status:      models.ScanStatusPending, // Set status to pending for further processing
		UserID:      userId,
		Sha256:      utils.SHA256(website.WebsiteUrl),
		SHA1:        utils.SHA1(website.WebsiteUrl),
		MD5:         utils.MD5(website.WebsiteUrl),
	}

	scanModel, err = service.CreateScan(scanModel)
	if err != nil {
		logger.Log.Error("Failed to save initial scan result: %v", err)
		return models.ScanModel{}, err
	}

	// Start background goroutines to handle further scans asynchronously
	go runBackgroundModules(scanModel.ID, url, website.Scripts)

	// Return the response immediately, without waiting for the background tasks
	return scanModel, nil
}

func runBackgroundModules(scanId uint, url string, scripts []models.ScriptRequest) {
	logger.Log.Info("Running background modules for scan ID: %d", scanId)

	go func() {
		logger.Log.Info("Starting network module for scan ID: %d", scanId)
		network.NetworkModule(scanId, url)
		logger.Log.Info("Completed network module for scan ID: %d", scanId)
	}()

	go func() {
		logger.Log.Info("Starting finding module for scan ID %d", scanId)
		finding.FindingModule(scanId, scripts)
		logger.Log.Info("Completed finding module for scan ID %d", scanId)
	}()

	go func() {
		logger.Log.Info("Starting content module for scan ID %d", scanId)
		content.ContentModule(scanId, scripts)
		logger.Log.Info("Completed content module for scan ID %d", scanId)
	}()

	go func() {
		logger.Log.Info("Starting list module for scan ID %d", scanId)
		list.ListModule(scanId, url)
		logger.Log.Info("Completed list module for scan ID %d", scanId)
	}()

	scanModel := models.ScanModel{
		Status: models.ScanStatusComplete,
	}
	scanModel.ID = scanId

	_, err := service.UpdateScan(scanModel)
	if err != nil {
		logger.Log.Error("Failed to update scan status: %v", err)
	}
}

// Helper function to format HTTP headers
func formatHTTPHeaders(headers http.Header) []string {
	httpHeaders := make([]string, 0)
	for key, value := range headers {
		httpHeaders = append(httpHeaders, fmt.Sprintf("%s: %s", key, value))
	}
	return httpHeaders
}

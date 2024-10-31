package core

import (
	"fmt"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/scanner"
)

type Inspector struct {
	scanRepository *repositories.ScanRepository
}

// NewInspector initializes the Inspector with necessary dependencies
func NewInspector(scanRepository *repositories.ScanRepository) *Inspector {
	return &Inspector{scanRepository: scanRepository}
}

// Execute performs the scan based on the job type
func (i *Inspector) Execute(job *entities.Job) error {
	switch job.Type {
	case "WebsiteScan":
		return i.performWebsiteScan(job)
	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}

// performWebsiteScan handles website scanning logic
func (i *Inspector) performWebsiteScan(job *entities.Job) error {
	// Perform the website scan
	scanResult, err := scanner.ScanWebsite(job.URL, job.UserAgent, uint(job.LicenseID))
	if err != nil {
		return err
	}

	// Save the scan result in DataStore via the ScanRepository
	return i.scanRepository.SaveScanResult(job.ID, scanResult)
}

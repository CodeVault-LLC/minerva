package core

import (
	"fmt"

	"github.com/codevault-llc/humblebrag-api/internal/core/modules"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"go.uber.org/zap"
)

type Inspector struct {
	modules map[string]modules.ScanModule
}

var InspectorCore *Inspector

// NewInspector initializes the Inspector with necessary dependencies
func NewInspector() *Inspector {
	inspector := &Inspector{modules: make(map[string]modules.ScanModule)}
	inspector.modules["network"] = modules.NewNetworkModule()
	inspector.modules["content"] = modules.NewContentModule()
	return inspector
}

// Execute performs the scan based on the job type
func (i *Inspector) Execute(job *entities.JobModel) error {
	switch job.Type {
	case "WebsiteScan":
		return i.performWebsiteScan(job)
	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}

// performWebsiteScan handles website scanning logic
func (i *Inspector) performWebsiteScan(job *entities.JobModel) error {
	logger.Log.Info("Starting website scan for URL: %s", zap.String("url", job.URL))

	requestedWebsite, err := websites.FetchWebsite(job.URL, job.UserAgent)
	if err != nil {
		return err
	}

	website, err := websites.AnalyzeHTML(requestedWebsite)
	if err != nil {
		logger.Log.Error("Failed to analyze website: %v", zap.Error(err))
		return err
	}

	scanModel := entities.ScanModel{
		Url:           job.URL,
		Title:         website.Title,
		RedirectChain: requestedWebsite.Redirects,
		StatusCode:    website.StatusCode,
		Status:        entities.ScanStatusPending,
		LicenseID:     uint(job.LicenseID),
		Sha256:        utils.SHA256(website.Url),
		SHA1:          utils.SHA1(website.Url),
		MD5:           utils.MD5(website.Url),
	}

	err = repository.ScanRepository.SaveScanResult(job, scanModel)
	if err != nil {
		logger.Log.Error("Failed to save scan result", zap.Error(err))
		return err
	}

	go func() {
		for _, module := range i.modules {
			if err := module.Execute(*job, website); err != nil {
				logger.Log.Error("Module failed", zap.Error(err), zap.String("module", module.Name()))
				return
			}
		}
	}()

	return nil
}

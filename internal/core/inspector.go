package core

import (
	"fmt"

	"github.com/codevault-llc/humblebrag-api/internal/core/modules"
	"github.com/codevault-llc/humblebrag-api/internal/database/storage"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/repository"
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
	inspector.modules["metadata"] = modules.NewMetadataModule()
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
	requestedWebsite, err := FetchWebsite(job.URL, job.UserAgent)
	if err != nil {
		return err
	}

	website, err := AnalyzeHTML(requestedWebsite)
	if err != nil {
		logger.Log.Error("Failed to analyze website: %v", zap.Error(err))
		return err
	}

	scanModel := entities.ScanModel{
		Url:        job.URL,
		Title:      website.Title,
		StatusCode: website.StatusCode,
		Status:     entities.ScanStatusPending,
		LicenseID:  uint(job.LicenseID),
		Sha256:     utils.SHA256(website.Url),
		SHA1:       utils.SHA1(website.Url),
		MD5:        utils.MD5(website.Url),
	}

	scan, err := repository.ScanRepository.SaveScanResult(job, scanModel)
	if err != nil {
		logger.Log.Error("Failed to save scan result", zap.Error(err))
		return err
	}

	job.ScanID = scan.ID

	for _, file := range website.Redirects {
		redirectModel := entities.RedirectModel{
			Url:        file.Url,
			HttpStatus: file.StatusCode,
			Timestamp:  utils.GetCurrentTime(),
			ScanID:     scan.ID,
		}

		redirect, err := repository.RedirectRepository.Create(redirectModel)
		if err != nil {
			logger.Log.Error("Failed to save redirect", zap.Error(err))
			return err
		}

		hashedBody := utils.SHA256(file.Screenshot.Content)
		err = storage.UploadFile("screenshot-bucket", hashedBody, []byte(file.Screenshot.Content), true)
		if err != nil {
			logger.Log.Error("Failed to upload screenshot", zap.Error(err))
			return err
		}

		_, err = repository.ScreenshotRepository.Create(entities.ScreenshotModel{
			RedirectId:     redirect.ID,
			ImageBucket:    "screenshot-bucket",
			ImageObjectKey: hashedBody,
			CompressedSize: len(file.Screenshot.Content),
		})
		if err != nil {
			logger.Log.Error("Failed to save screenshot", zap.Error(err))
			return err
		}
	}

	go func() {
		for _, module := range i.modules {
			if err := module.Execute(*job, website); err != nil {
				logger.Log.Error("Module failed", zap.Error(err), zap.String("module", module.Name()))
				continue
			}
		}
	}()

	return nil
}

package contents

import (
	"sync"
	"time"

	"github.com/codevault-llc/minerva/config"
	"github.com/codevault-llc/minerva/internal/common"
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	repository "github.com/codevault-llc/minerva/internal/contents/models/repository"
	generalEntities "github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/database/storage"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContentModule struct {
	runtimeLocation common.RuntimeLocation
	repository      *repository.ContentRepo
	findingRepo     *repository.FindingRepo
}

func NewContentModule(runtimeLocation common.RuntimeLocation, db *database.Database) *ContentModule {
	repository.ContentRepository = repository.NewContentRepo(db)
	repository.FindingRepository = repository.NewFindingRepo(db)

	return &ContentModule{
		runtimeLocation: runtimeLocation,
		repository:      repository.ContentRepository,
		findingRepo:     repository.FindingRepository,
	}
}

func (m *ContentModule) Execute(job generalEntities.JobModel, website *common.WebsiteAnalysis) error {
	for _, script := range website.Assets {
		hashedBody := utils.SHA256(script.Content)

		existingContentId, err := m.repository.FindContentByHash(hashedBody)
		if err != nil {
			logger.Log.Error("Failed to find content by hash: %v", zap.Error(err))

			var jsFiles []common.FileRequest
			for _, asset := range website.Assets {
				if asset.FileType == "application/javascript" {
					jsFiles = append(jsFiles, asset)
				}
			}

			findings := scanSecrets(jsFiles)
			err := m.findingRepo.SaveFindingResult(job, findings)
			if err != nil {
				logger.Log.Error("Failed to save finding result: %v", zap.Error(err))
			}

			continue
		}

		if existingContentId == "" {
			originalFileName := script.Src
			fileExtension := storage.GetFileExtension(originalFileName)
			objectKey := storage.GenerateObjectKey(originalFileName)
			sanitizedObjectKey := storage.SanitizeObjectKey(objectKey)
			contentType := storage.GetContentType(fileExtension)

			err = storage.UploadFile("content-bucket", sanitizedObjectKey, []byte(script.Content), contentType, true)
			if err != nil {
				logger.Log.Error("Failed to upload file: %v", zap.Error(err))
				continue
			}

			content := entities.ContentModel{
				Id:          uuid.New().String(),
				ScanId:      job.ScanID,
				FileSize:    int64(script.FileSize),
				FileType:    script.FileType,
				Source:      script.Src,
				StorageType: storage.DetermineStorageType(script.Content),
				HashedBody:  hashedBody,
				Tags:        []string{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			err := m.repository.SaveContentResult(content)
			if err != nil {
				logger.Log.Error("Failed to save content: %v", zap.Error(err))
				continue
			}

			storageRecord := entities.ContentStorageModel{
				Id:              uuid.New().String(),
				ContentId:       content.Id,
				BucketName:      "content-bucket",
				ObjectKey:       sanitizedObjectKey,
				Location:        storage.GetLocation("content-bucket", sanitizedObjectKey),
				StorageEndpoint: storage.GetEndpoint("content-bucket"),
				Encryption:      "AES256",
			}

			err = m.repository.CreateContentStorage(storageRecord)
			if err != nil {
				logger.Log.Error("Failed to save storage record: %v", zap.Error(err))
				continue
			}
		}
	}

	var jsFiles []common.FileRequest
	for _, asset := range website.Assets {
		if asset.FileType == string(utils.ApplicationJavascript) {
			jsFiles = append(jsFiles, asset)
		}
	}

	findings := scanSecrets(jsFiles)
	err := m.findingRepo.SaveFindingResult(job, findings)
	if err != nil {
		logger.Log.Error("Failed to save finding result: %v", zap.Error(err))
	}

	return nil
}

func (m *ContentModule) Name() string {
	return "content"
}

func (m *ContentModule) RuntimeLocation() common.RuntimeLocation {
	return m.runtimeLocation
}

func scanSecrets(scripts []common.FileRequest) []utils.RegexReturn {
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
				matches := utils.GenericScan(rule, script)
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

package contents

import (
	"sync"
	"time"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/contents/models/entities"
	repository "github.com/codevault-llc/humblebrag-api/internal/contents/models/repository"
	generalEntities "github.com/codevault-llc/humblebrag-api/internal/core/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/database/storage"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"go.uber.org/zap"
)

type ContentModule struct{}

func NewContentModule() *ContentModule {
	return &ContentModule{}
}

func (m *ContentModule) Execute(job generalEntities.JobModel, website types.WebsiteAnalysis) error {
	for _, script := range website.Assets {
		hashedBody := utils.SHA256(script.Content)

		existingContent, err := repository.ContentRepository.FindContentByHash(hashedBody)
		if err != nil {
			logger.Log.Error("Failed to find content by hash: %v", zap.Error(err))

			var jsFiles []types.FileRequest
			for _, asset := range website.Assets {
				if asset.FileType == "application/javascript" {
					jsFiles = append(jsFiles, asset)
				}
			}

			findings := scanSecrets(jsFiles)
			repository.FindingRepository.SaveFindingResult(job, findings)
			continue
		}

		if existingContent.Id != 0 {
			err := repository.ContentRepository.IncrementAccessCount(existingContent.Id)
			if err != nil {
				logger.Log.Error("Failed to increment access count: %v", zap.Error(err))
			}
		} else {
			storageType := storage.DetermineStorageType(script.Content)
			err = storage.UploadFile("content-bucket", hashedBody, []byte(script.Content), true)
			if err != nil {
				logger.Log.Error("Failed to upload file: %v", zap.Error(err))
				continue
			}

			content := entities.ContentModel{
				ScanId:         job.ScanID,
				FileSize:       int64(script.FileSize),
				FileType:       script.FileType,
				Source:         script.Src,
				StorageType:    storageType,
				AccessCount:    1,
				HashedBody:     hashedBody,
				LastAccessedAt: time.Now(),
			}

			newContent, err := repository.ContentRepository.SaveContentResult(content)
			if err != nil {
				logger.Log.Error("Failed to save content: %v", zap.Error(err))
				continue
			}

			storageRecord := entities.ContentStorageModel{
				ContentId:       newContent,
				BucketName:      "content-bucket",
				ObjectKey:       hashedBody,
				Location:        storage.GetLocation("content-bucket", hashedBody),
				StorageEndpoint: storage.GetEndpoint("content-bucket"),
				Encryption:      "AES256",
			}

			err = repository.ContentRepository.CreateContentStorage(storageRecord)
			if err != nil {
				logger.Log.Error("Failed to save storage record: %v", zap.Error(err))
				continue
			}
		}
	}

	var jsFiles []types.FileRequest
	for _, asset := range website.Assets {
		if asset.FileType == "application/javascript" {
			jsFiles = append(jsFiles, asset)
		}
	}

	findings := scanSecrets(jsFiles)
	repository.FindingRepository.SaveFindingResult(job, findings)
	return nil
}

func (m *ContentModule) Name() string {
	return "content"
}

func scanSecrets(scripts []types.FileRequest) []utils.RegexReturn {
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

package content

import (
	"sync"
	"time"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/database/storage"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"go.uber.org/zap"
)

func ContentModule(scanId uint, assets []types.FileRequest) {
	for _, script := range assets {
		hashedBody := utils.SHA256(script.Content)

		// Check if content already exists by hash
		existingContent, err := service.FindContentByHash(hashedBody)
		if err != nil {
			logger.Log.Error("Failed to find content by hash: %v", zap.Error(err))
			continue
		}

		var contentID uint

		if existingContent.ID != 0 {
			// Increment access count if the content already exists
			err := service.IncrementAccessCount(existingContent.ID)
			if err != nil {
				logger.Log.Error("Failed to increment access count: %v", zap.Error(err))
			}
			contentID = existingContent.ID
		} else {
			// Upload the new file to S3
			storageType := storage.DetermineStorageType(script.Content)
			err = storage.UploadFile("content-bucket", hashedBody, []byte(script.Content))
			if err != nil {
				logger.Log.Error("Failed to upload file: %v", zap.Error(err))
				continue
			}

			// Create new content record
			content := models.ContentModel{
				FileSize:       int64(script.FileSize),
				FileType:       script.FileType,
				StorageType:    storageType,
				AccessCount:    1,
				HashedBody:     hashedBody,
				LastAccessedAt: time.Now(),
			}

			newContent, err := service.CreateContent(content)
			if err != nil {
				logger.Log.Error("Failed to save content: %v", zap.Error(err))
				continue
			}

			storageRecord := models.ContentStorageModel{
				ContentID:       newContent.ID,
				BucketName:      "content-bucket",
				ObjectKey:       hashedBody,
				Location:        storage.GetLocation("content-bucket", hashedBody),
				StorageEndpoint: storage.GetEndpoint("content-bucket"),
				Encryption:      "AES256",
			}

			err = service.CreateContentStorage(storageRecord)
			if err != nil {
				logger.Log.Error("Failed to save storage record: %v", zap.Error(err))
				continue
			}

			contentID = newContent.ID
		}

		// Associate the content with the scan
		err = service.AddContentToScan(scanId, contentID)
		if err != nil {
			logger.Log.Error("Failed to associate content with scan: %v", zap.Error(err))
		}
	}

	findings := scanSecrets(assets)
	service.CreateFindings(scanId, findings)
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

package contents

import (
	"context"
	"sync"

	"github.com/codevault-llc/minerva/config"
	"github.com/codevault-llc/minerva/internal/common"
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	repository "github.com/codevault-llc/minerva/internal/contents/models/repository"
	generalEntities "github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/fingerprint"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"

	pb "github.com/codevault-llc/minerva/proto"
)

type ContentModule struct {
	runtimeLocation common.RuntimeLocation
	repository      *repository.ContentRepo
	findingRepo     *repository.FindingRepo
	fingerprintRepo *repository.FingerprintRepo
}

func NewContentModule(runtimeLocation common.RuntimeLocation, db *database.Database) *ContentModule {
	repository.ContentRepository = repository.NewContentRepo(db)
	repository.FindingRepository = repository.NewFindingRepo(db)
	repository.FingerprintRepository = repository.NewFingerprintRepo(db)

	return &ContentModule{
		runtimeLocation: runtimeLocation,
		repository:      repository.ContentRepository,
		findingRepo:     repository.FindingRepository,
		fingerprintRepo: repository.FingerprintRepository,
	}
}

func (m *ContentModule) Execute(job generalEntities.JobModel, website *common.WebsiteAnalysis) error {
	for _, script := range website.Assets {
		md5Data := utils.MD5(script.Src)

		existingContentId, err := m.repository.FindContentByMd5(md5Data)
		if err != nil {
			logger.Log.Error("Failed to find content by hash: %v", zap.Error(err))

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

			continue
		}

		if existingContentId == "" {
			existingContentId = uuid.New().String()

			content := entities.ContentModel{
				Id:       existingContentId,
				ScanId:   job.ScanID,
				FileSize: int64(script.FileSize),
				FileType: script.FileType,
				Source:   script.Src,
				Md5:      md5Data,
				Sha1:     utils.SHA1(script.Src),
				Sha256:   utils.SHA256(script.Src),
				Tags:     []string{},
			}

			err := m.repository.SaveContentResult(content)
			if err != nil {
				logger.Log.Error("Failed to save content: %v", zap.Error(err))
				continue
			}
		}

		foundFingerprints, err := fingerprint.FingerprintClient.MatchFingerprint(context.Background(), &pb.MatchFingerprintRequest{
			Source: script.Src,
		})

		if err != nil {
			logger.Log.Error("Failed to match fingerprint: %v", zap.Error(err))
			continue
		}

		for _, fingerprint := range foundFingerprints.Matched {
			logger.Log.Info("Matched fingerprint", zap.String("fingerprint_id", fingerprint.Id))

			fingerprintRecord := entities.FingerprintModel{
				ContentId:     existingContentId,
				FingerprintId: fingerprint.Id,
			}

			err = m.fingerprintRepo.SaveFingerprintResult(fingerprintRecord)
			if err != nil {
				logger.Log.Error("Failed to save fingerprint result: %v", zap.Error(err))
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

	for _, rule := range config.Config.Rules {
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
		}(rule)
	}

	wg.Wait()
	return results
}

package repository

import (
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FingerprintRepo struct {
	database *database.Database
}

func NewFingerprintRepo(database *database.Database) *FingerprintRepo {
	return &FingerprintRepo{database: database}
}

var FingerprintRepository *FingerprintRepo

func (repository *FingerprintRepo) SaveFingerprintResult(contentId string, fingerprints []utils.RegexReturn) error {
	for _, finding := range fingerprints {
		for _, match := range finding.Matches {
			finding := entities.FingerprintModel{
				Id:        uuid.New().String(),
				ContentId: contentId,

				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				FingerprintName: finding.Name,

				CreatedAt: utils.GetCurrentTime(),
				UpdatedAt: utils.GetCurrentTime(),
			}

			query := "INSERT INTO fingerprint (id, content_id, line, match, source, fingerprint_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
			queryResult := repository.database.GetDatabase().Query(query, finding.Id, finding.ContentId, finding.Line, finding.Match, finding.Source, finding.FingerprintName, finding.CreatedAt, finding.UpdatedAt)
			err := queryResult.Exec()
			if err != nil {
				logger.Log.Error("Failed to save finding", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

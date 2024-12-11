package repository

import (
	"context"
	"errors"

	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	generalEntities "github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FindingRepo struct {
	database *database.Database
}

func NewFindingRepo(database *database.Database) *FindingRepo {
	return &FindingRepo{database: database}
}

var FindingRepository *FindingRepo

func (repository *FindingRepo) SaveFindingResult(job generalEntities.JobModel, findings []utils.RegexReturn) error {
	for _, finding := range findings {
		for _, match := range finding.Matches {
			finding := entities.FindingModel{
				Id:     uuid.New().String(),
				ScanId: job.ScanID,

				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				RegexName:        finding.Name,
				RegexDescription: finding.Description,

				CreatedAt: utils.GetCurrentTime(),
				UpdatedAt: utils.GetCurrentTime(),
			}

			query := "INSERT INTO finding (id, scan_id, line, match, source, regex_name, regex_description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
			queryResult := repository.database.GetDatabase().Query(query, finding.Id, finding.ScanId, finding.Line, finding.Match, finding.Source, finding.RegexName, finding.RegexDescription, finding.CreatedAt, finding.UpdatedAt)
			err := queryResult.Exec()
			if err != nil {
				logger.Log.Error("Failed to save finding", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

func (repository *FindingRepo) GetScanFindings(scanID uint) ([]entities.FindingModel, error) {
	ctx := context.Background()
	var findings []entities.FindingModel

	query := "SELECT id, scan_id, line, match, source, regex_name, regex_description, created_at, updated_at FROM finding WHERE scan_id = ?"
	err := repository.database.Select(ctx, query, &findings, scanID)
	if err != nil {
		logger.Log.Error("Failed to fetch scan result", zap.Error(err))
		return []entities.FindingModel{}, err
	}

	if len(findings) == 0 {
		return []entities.FindingModel{}, errors.New("no finding result found")
	}

	return findings, nil
}

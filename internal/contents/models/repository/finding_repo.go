package repository

import (
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	generalEntities "github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
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
				ScanId: job.ScanID,

				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				RegexName:        finding.Name,
				RegexDescription: finding.Description,
			}

			query := "INSERT INTO finding (scan_id, line, match, source, regex_name, regex_description) VALUES ($1, $2, $3, $4, $5, $6)"
			queryResult := repository.database.GetDatabase().Query(query, finding.ScanId, finding.Line, finding.Match, finding.Source, finding.RegexName, finding.RegexDescription)
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
	var findings []entities.FindingModel

	query := "SELECT * FROM finding WHERE scan_id = $1"
	queryResult := repository.database.GetDatabase().Query(query, scanID)

	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to get scan findings", zap.Error(err))
		return nil, err
	}

	err = queryResult.Scan(&findings)
	if err != nil {
		logger.Log.Error("Failed to scan scan findings", zap.Error(err))
		return nil, err
	}

	return findings, nil
}

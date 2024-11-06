package repository

import (
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	generalEntities "github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type FindingRepo struct {
	db *sqlx.DB
}

func NewFindingRepo(db *sqlx.DB) *FindingRepo {
	return &FindingRepo{db: db}
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

			tx, err := repository.db.Beginx()
			if err != nil {
				return err
			}

			query, values, err := database.StructToQuery(finding, "finding")
			if err != nil {
				return err
			}

			_, err = database.InsertStruct(tx, query, values)
			if err != nil {
				return err
			}

			tx.Commit()
		}
	}

	return nil
}

func (repository *FindingRepo) GetScanFindings(scanID uint) ([]entities.FindingModel, error) {
	var findings []entities.FindingModel

	err := repository.db.Select(&findings, "SELECT * FROM finding WHERE scan_id = $1", scanID)
	if err != nil {
		logger.Log.Error("Failed to get scan findings", zap.Error(err))
		return nil, err
	}

	return findings, nil
}

package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/contents/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/database"
	generalEntities "github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/jmoiron/sqlx"
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

			query, err := database.StructToQuery(finding, "finding")
			if err != nil {
				return err
			}

			_, err = database.InsertStruct(tx, query, finding)
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

	repository.db.Get(&findings, "SELECT * FROM finding WHERE scan_id = $1", scanID)

	return findings, nil
}

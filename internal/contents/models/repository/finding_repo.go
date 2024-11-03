package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/contents/models/entities"
	generalEntities "github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"gorm.io/gorm"
)

type FindingRepo struct {
	db *gorm.DB
}

func NewFindingRepo(db *gorm.DB) *FindingRepo {
	return &FindingRepo{db: db}
}

var FindingRepository *FindingRepo

func (repository *FindingRepo) SaveFindingResult(job generalEntities.JobModel, findings []utils.RegexReturn) error {
	for _, finding := range findings {
		for _, match := range finding.Matches {
			finding := entities.FindingModel{
				ScanID: job.ScanID,

				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				RegexName:        finding.Name,
				RegexDescription: finding.Description,
			}

			tx := repository.db.Begin()
			if err := tx.Create(&finding).Error; err != nil {
				tx.Rollback()
				return err
			}

			tx.Commit()
		}
	}

	return nil
}

func (repository *FindingRepo) GetScanFindings(scanID uint) ([]entities.FindingModel, error) {
	var findings []entities.FindingModel

	if err := repository.db.Where("scan_id = ?", scanID).Find(&findings).Error; err != nil {
		return nil, err
	}

	return findings, nil
}

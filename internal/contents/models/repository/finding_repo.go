package repository

import (
	"context"
	"errors"
	"time"

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
	// Prepare the batch insertion query
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `
		INSERT INTO finding (id, scan_id, line, match, source, regex_name, regex_description)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Prepare a batch for efficient insertion
	batch, err := repository.database.GetDatabase().PrepareBatch(ctx, query)
	if err != nil {
		logger.Log.Error("Failed to prepare batch for findings", zap.Error(err))
		return err
	}

	for _, finding := range findings {
		for _, match := range finding.Matches {
			newFinding := entities.FindingModel{
				Id:               uuid.New().String(),
				ScanId:           job.ScanID,
				Line:             match.Line,
				Match:            match.Match,
				Source:           match.Source,
				RegexName:        finding.Name,
				RegexDescription: finding.Description,
			}

			// Add to batch
			err := batch.Append(
				newFinding.Id,
				newFinding.ScanId,
				newFinding.Line,
				newFinding.Match,
				newFinding.Source,
				newFinding.RegexName,
				newFinding.RegexDescription,
			)
			if err != nil {
				logger.Log.Error("Failed to append finding to batch", zap.Error(err))
				return err
			}
		}
	}

	// Execute the batch
	if err := batch.Flush(); err != nil {
		logger.Log.Error("Failed to execute batch for findings", zap.Error(err))
		return err
	}

	return nil
}

func (repository *FindingRepo) GetScanFindings(scanID uint) ([]entities.FindingModel, error) {
	ctx := context.Background()
	var findings []entities.FindingModel

	query := "SELECT id, scan_id, line, match, source, regex_name, regex_description, created_at, updated_at FROM finding WHERE scan_id = ?"
	err := repository.database.Db.QueryRow(ctx, query, scanID).Scan(&findings)
	if err != nil {
		logger.Log.Error("Failed to fetch scan result", zap.Error(err))
		return []entities.FindingModel{}, err
	}

	if len(findings) == 0 {
		return []entities.FindingModel{}, errors.New("no finding result found")
	}

	return findings, nil
}

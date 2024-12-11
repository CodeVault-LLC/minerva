package repository

import (
	"context"
	"errors"

	"github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"go.uber.org/zap"
)

type ScanRepo struct {
	database *database.Database
}

// NewScanRepository creates a new ScanRepository
func NewScanRepository(database *database.Database) *ScanRepo {
	return &ScanRepo{
		database: database,
	}
}

var ScanRepository *ScanRepo

// SaveScanResult saves the scan result in the database
func (repository *ScanRepo) SaveScanResult(job *entities.JobModel, scan entities.ScanModel) error {
	query := "INSERT INTO minerva.scans (id, url, title, status_code, status, sha256, sha1, md5, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	queryResult := repository.database.GetDatabase().Query(query, scan.Id, job.URL, scan.Title, scan.StatusCode, scan.Status, scan.Sha256, scan.Sha1, scan.Md5, scan.CreatedAt, scan.UpdatedAt)
	err := queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return err
	}

	return nil
}

func (repository *ScanRepo) GetScanResult(scanId string) (entities.ScanModel, error) {
	ctx := context.Background()
	var scans []entities.ScanModel

	query := "SELECT id, url, title, status_code, status, sha256, sha1, md5, created_at, updated_at FROM scans WHERE id = ?"
	err := repository.database.Select(ctx, query, &scans, scanId)
	if err != nil {
		logger.Log.Error("Failed to fetch scan result", zap.Error(err))
		return entities.ScanModel{}, err
	}

	if len(scans) == 0 {
		return entities.ScanModel{}, errors.New("no scan result found")
	}

	return scans[0], nil
}

func (repository *ScanRepo) GetScans() ([]entities.ScanModel, error) {
	ctx := context.Background()
	var scans []entities.ScanModel

	query := "SELECT id, url, title, status_code, status, sha256, sha1, md5, created_at, updated_at FROM scans"
	err := repository.database.Select(ctx, query, &scans)
	if err != nil {
		logger.Log.Error("Failed to fetch scan result", zap.Error(err))
		return []entities.ScanModel{}, err
	}

	if len(scans) == 0 {
		return []entities.ScanModel{}, errors.New("no scan result found")
	}

	return scans, nil
}

func (repository *ScanRepo) CompleteScan(scanId string) error {
	query := "UPDATE scans SET status = ? WHERE id = ?"
	queryResult := repository.database.GetDatabase().Query(query, entities.ScanStatusComplete, scanId)
	err := queryResult.Exec()
	if err != nil {
		return err
	}

	return nil
}

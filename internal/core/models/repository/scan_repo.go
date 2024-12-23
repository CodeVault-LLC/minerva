package repository

import (
	"context"
	"errors"
	"time"

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
	query := `
		INSERT INTO scans (id, url, title, status_code, sha256, sha1, md5)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute the insertion
	err := repository.database.GetDatabase().Exec(ctx, query,
		job.ScanID,
		scan.Url,
		scan.Title,
		scan.StatusCode,
		scan.Sha256,
		scan.Sha1,
		scan.Md5,
	)
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
	err := repository.database.Db.QueryRow(ctx, query, scanId).Scan(&scans)
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
	err := repository.database.Db.QueryRow(ctx, query).Scan(&scans)
	if err != nil {
		logger.Log.Error("Failed to fetch scan result", zap.Error(err))
		return []entities.ScanModel{}, err
	}

	if len(scans) == 0 {
		return []entities.ScanModel{}, errors.New("no scan result found")
	}

	return scans, nil
}

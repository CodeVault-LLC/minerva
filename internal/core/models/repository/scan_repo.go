package repository

import (
	"github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ScanRepo struct {
	db *sqlx.DB
}

// NewScanRepository creates a new ScanRepository
func NewScanRepository(db *sqlx.DB) *ScanRepo {
	return &ScanRepo{
		db: db,
	}
}

var ScanRepository *ScanRepo

// SaveScanResult saves the scan result in the database
func (repository *ScanRepo) SaveScanResult(job *entities.JobModel, scan entities.ScanModel) (uint, error) {
	tx, err := repository.db.Beginx()
	if err != nil {
		return 0, err
	}

	query, values, err := database.StructToQuery(scan, "scans")
	if err != nil {
		logger.Log.Error("Failed to generate query", zap.Error(err))
		return 0, err
	}

	returnId, err := database.InsertStruct(tx, query, values)
	if err != nil {
		logger.Log.Error("Failed to insert certificate", zap.Error(err))
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error("Failed to commit transaction", zap.Error(err))
		return 0, err
	}

	return returnId, nil
}

func (repository *ScanRepo) GetScanResult(scanId uint) (entities.ScanModel, error) {
	var scan entities.ScanModel
	err := repository.db.Get(&scan, "SELECT * FROM scans WHERE id = $1", scanId)
	if err != nil {
		return scan, err
	}

	return scan, nil
}

func (repository *ScanRepo) GetScans() ([]entities.ScanModel, error) {
	var scans []entities.ScanModel
	err := repository.db.Select(&scans, "SELECT * FROM scans")
	if err != nil {
		logger.Log.Error("Failed to get scans", zap.Error(err))
		return scans, err
	}

	return scans, nil
}

func (repository *ScanRepo) CompleteScan(scanId uint) error {
	_, err := repository.db.Exec("UPDATE scans SET status = $1 WHERE id = $2", entities.ScanStatusComplete, scanId)
	if err != nil {
		return err
	}

	return nil
}

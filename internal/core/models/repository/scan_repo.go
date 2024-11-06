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
	query := "SELECT * FROM scans WHERE id = $1"
	stmt, err := repository.db.Preparex(query)
	if err != nil {
		return entities.ScanModel{}, err
	}

	var scan entities.ScanModel
	err = stmt.Get(&scan, scanId)
	if err != nil {
		logger.Log.Error("Failed to get scan", zap.Error(err))
		return entities.ScanModel{}, err
	}

	return scan, nil
}

func (repository *ScanRepo) GetScans() ([]entities.ScanModel, error) {
	query := "SELECT * FROM scans"
	stmt, err := repository.db.Preparex(query)
	if err != nil {
		return nil, err
	}

	var scans []entities.ScanModel
	err = stmt.Select(&scans)
	if err != nil {
		logger.Log.Error("Failed to get scan", zap.Error(err))
		return nil, err
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

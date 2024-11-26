package repository

import (
	"github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
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
func (repository *ScanRepo) SaveScanResult(job *entities.JobModel, scan entities.ScanModel) (uint, error) {
	query := "INSERT INTO scans (job_id, status, url) VALUES (?, ?, ?) RETURNING id"

	queryResult := repository.database.GetDatabase().Query(query, job.ID, entities.ScanStatusComplete, scan.Url)
	err := queryResult.Exec()
	if err != nil {
		return 0, err
	}

	var scanId uint
	err = queryResult.Scan(&scanId)
	if err != nil {
		return 0, err
	}

	return scanId, nil
}

func (repository *ScanRepo) GetScanResult(scanId uint) (entities.ScanModel, error) {
	query := "SELECT * FROM scans WHERE id = ?"
	queryResult := repository.database.GetDatabase().Query(query, scanId)

	var scan entities.ScanModel
	err := queryResult.Scan(&scan)
	if err != nil {
		return entities.ScanModel{}, err
	}

	return scan, nil
}

func (repository *ScanRepo) GetScans() ([]entities.ScanModel, error) {
	query := "SELECT * FROM scans"
	queryResult := repository.database.GetDatabase().Query(query)

	var scans []entities.ScanModel
	err := queryResult.Scan(&scans)
	if err != nil {
		return nil, err
	}

	return scans, nil
}

func (repository *ScanRepo) CompleteScan(scanId uint) error {
	query := "UPDATE scans SET status = ? WHERE id = ?"
	queryResult := repository.database.GetDatabase().Query(query, entities.ScanStatusComplete, scanId)
	err := queryResult.Exec()
	if err != nil {
		return err
	}

	return nil
}

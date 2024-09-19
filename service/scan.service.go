package service

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func CreateScan(scan models.ScanModel) (models.ScanModel, error) {
	constants.DB.Model(&models.ScanModel{}).Where("website_url = ?", scan.WebsiteUrl).Update("status", models.ScanStatusArchived)

	if err := constants.DB.Create(&scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

func GetScans() ([]models.ScanAPIResponse, error) {
	var scans []models.ScanModel

	if err := constants.DB.Preload("Findings").Where("status = ?", models.ScanStatusComplete).Order("created_at desc").Find(&scans).Error; err != nil {
		return models.ConvertScans(scans), err
	}

	return models.ConvertScans(scans), nil
}

func GetScan(scanID string) (models.ScanAPIResponse, error) {
	var scan models.ScanModel

	if err := constants.DB.Where("id = ?", scanID).Preload("Findings").Preload("Lists").Preload("Certificates").Where("status = ?", models.ScanStatusComplete).First(&scan).Error; err != nil {
		return models.ConvertScan(scan), err
	}

	return models.ConvertScan(scan), nil
}

func GetScansByUserID(userID uint) ([]models.ScanModel, error) {
	var scans []models.ScanModel

	if err := constants.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&scans).
		Error; err != nil {
		return scans, err
	}

	return scans, nil
}

// GetTotalScans returns the total number of scans
func GetTotalScans() (int64, error) {
	var totalScans int64
	if err := constants.DB.Model(&models.ScanModel{}).Count(&totalScans).Error; err != nil {
		return 0, err
	}

	return totalScans, nil
}

// GetTotalDomainsScanned returns the total number of unique domains scanned
func GetTotalDomainsScanned() (int64, error) {
	var totalDomainsScanned int64
	if err := constants.DB.Model(&models.ScanModel{}).Distinct("website_url").Count(&totalDomainsScanned).Error; err != nil {
		return 0, err
	}

	return totalDomainsScanned, nil
}

// GetRecentScans returns the number of scans in the last 24 hours
func GetRecentScans() (int64, error) {
	var lastScansIn24Hours int64
	if err := constants.DB.Model(&models.ScanModel{}).Where("created_at >= ?", utils.Get24HoursAgo()).Count(&lastScansIn24Hours).Error; err != nil {
		return 0, err
	}

	return lastScansIn24Hours, nil
}

// GetMostScannedDomains returns the top 10 most scanned domains
func GetMostScannedDomains() ([]models.ScanModel, error) {
	var scans []models.ScanModel
	if err := constants.DB.Model(&models.ScanModel{}).Select("*, count(website_url) as count").Group("website_url, id").Order("count desc").Limit(10).Find(&scans).Error; err != nil {
		return scans, err
	}

	return scans, nil
}

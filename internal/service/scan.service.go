package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func CreateScan(scan models.ScanModel) (models.ScanModel, error) {
	database.DB.Model(&models.ScanModel{}).Where("website_url = ?", scan.WebsiteUrl).Update("status", models.ScanStatusArchived)

	if err := database.DB.Create(&scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

// Make this function update the scan information. Now we are just updating some fields of the scan, meaning some can be empty.
func UpdateScan(scan models.ScanModel) (models.ScanModel, error) {
	if err := database.DB.Model(&models.ScanModel{}).Where("id = ?", scan.ID).Updates(scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

func GetScans() ([]models.ScanAPIResponse, error) {
	var scans []models.ScanModel

	if err := database.DB.Preload("Findings").Where("status IN (?, ?)", models.ScanStatusComplete, models.ScanStatusPending).Order("created_at desc").Find(&scans).Error; err != nil {
		return models.ConvertScans(scans), err
	}

	return models.ConvertScans(scans), nil
}

func GetScan(scanID string) (models.ScanAPIResponse, error) {
	var scan models.ScanModel

	if err := database.DB.Where("id = ?", scanID).Preload("Findings").Preload("Lists").Where("status IN (?, ?)", models.ScanStatusComplete, models.ScanStatusPending).First(&scan).Error; err != nil {
		return models.ConvertScan(scan), err
	}

	return models.ConvertScan(scan), nil
}

func GetScansByUserID(userID uint) ([]models.ScanModel, error) {
	var scans []models.ScanModel

	if err := database.DB.Where("user_id = ?", userID).
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
	if err := database.DB.Model(&models.ScanModel{}).Count(&totalScans).Error; err != nil {
		return 0, err
	}

	return totalScans, nil
}

// GetTotalDomainsScanned returns the total number of unique domains scanned
func GetTotalDomainsScanned() (int64, error) {
	var totalDomainsScanned int64
	if err := database.DB.Model(&models.ScanModel{}).Distinct("website_url").Count(&totalDomainsScanned).Error; err != nil {
		return 0, err
	}

	return totalDomainsScanned, nil
}

// GetRecentScans returns the number of scans in the last 24 hours
func GetRecentScans() (int64, error) {
	var lastScansIn24Hours int64
	if err := database.DB.Model(&models.ScanModel{}).Where("created_at >= ?", utils.Get24HoursAgo()).Count(&lastScansIn24Hours).Error; err != nil {
		return 0, err
	}

	return lastScansIn24Hours, nil
}

// GetMostScannedDomains returns the top 10 most scanned domains
func GetMostScannedDomains() ([]models.ScanModel, error) {
	var scans []models.ScanModel
	if err := database.DB.Model(&models.ScanModel{}).Select("*, count(website_url) as count").Group("website_url, id").Order("count desc").Limit(10).Find(&scans).Error; err != nil {
		return scans, err
	}

	return scans, nil
}

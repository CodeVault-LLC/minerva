package controller

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func CreateScan(scan models.Scan) (models.Scan, error) {
	if err := constants.DB.Create(&scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

func CreateFindings(scanID uint, secrets []utils.RegexReturn) {
	for _, secret := range secrets {
		for _, match := range secret.Matches {
			finding := models.Finding{
				ScanID: scanID,
				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				RegexName:        secret.Name,
				RegexDescription: secret.Description,
			}

			constants.DB.Create(&finding)
		}
	}
}

func CreateContent(content models.Content) (models.Content, error) {
	if err := constants.DB.Create(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

func GetScans() ([]models.ScanResponse, error) {
	var scans []models.Scan

	if err := constants.DB.Preload("User").Preload("Findings").Find(&scans).Error; err != nil {
		return utils.ConvertScans(scans), err
	}

	return utils.ConvertScans(scans), nil
}

func GetScan(scanID string) (models.ScanResponse, error) {
	var scan models.Scan

	if err := constants.DB.Where("id = ?", scanID).
		First(&scan).
		Error; err != nil {
		return utils.ConvertScan(scan), err
	}

	return utils.ConvertScan(scan), nil
}

func GetScanFindings(scanID string) ([]models.FindingResponse, error) {
	var findings []models.Finding

	if err := constants.DB.Where("scan_id = ?", scanID).
		Find(&findings).
		Error; err != nil {
		return utils.ConvertFindings(findings), err
	}

	return utils.ConvertFindings(findings), nil
}

func GetScanContent(scanID string) ([]models.ContentResponse, error) {
	var content []models.Content

	if err := constants.DB.Where("scan_id = ?", scanID).
		Find(&content).
		Error; err != nil {
		return utils.ConvertContents(content), err
	}

	return utils.ConvertContents(content), nil
}

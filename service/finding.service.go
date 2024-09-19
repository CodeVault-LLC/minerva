package service

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
)

// CreateFindings creates findings in the database
func CreateFindings(scanID uint, secrets []utils.RegexReturn) {
	for _, secret := range secrets {
		for _, match := range secret.Matches {
			finding := models.FindingModel{
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

// GetScanFindings retrieves findings from the database
func GetScanFindings(scanID string) ([]models.FindingResponse, error) {
	var findings []models.FindingModel

	if err := constants.DB.Where("scan_id = ?", scanID).
		Find(&findings).
		Error; err != nil {
		return models.ConvertFindings(findings), err
	}

	return models.ConvertFindings(findings), nil
}

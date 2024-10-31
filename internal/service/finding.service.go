package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

// CreateFindings creates findings in the database
func CreateFindings(scanID uint, secrets []utils.RegexReturn) {
	for _, secret := range secrets {
		for _, match := range secret.Matches {
			finding := entities.FindingModel{
				ScanID: scanID,
				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				RegexName:        secret.Name,
				RegexDescription: secret.Description,
			}

			database.DB.Create(&finding)
		}
	}
}

// UpdateFindings	updates findings in the database
func UpdateFindings(scanID uint, secrets []utils.RegexReturn) {
	var findings []entities.FindingModel

	database.DB.Where("scan_id = ?", scanID).Find(&findings)

	for _, secret := range secrets {
		for _, match := range secret.Matches {
			finding := entities.FindingModel{
				ScanID: scanID,
				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				RegexName:        secret.Name,
				RegexDescription: secret.Description,
			}

			if !viewmodels.FindFinding(findings, finding) {
				database.DB.Create(&finding)
			}
		}
	}
}

// GetScanFindings retrieves findings from the database
func GetScanFindings(scanID uint) ([]viewmodels.Finding, error) {
	var findings []entities.FindingModel

	if err := database.DB.Where("scan_id = ?", scanID).
		Find(&findings).
		Error; err != nil {
		return viewmodels.ConvertFindings(findings), err
	}

	return viewmodels.ConvertFindings(findings), nil
}

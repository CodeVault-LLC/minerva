package models

import "gorm.io/gorm"

type FindingModel struct {
	gorm.Model

	ScanID uint
	Scan   ScanModel

	RegexName        string `gorm:"not null"`
	RegexDescription string `gorm:"not null"`

	Match  string `gorm:"not null"`
	Source string `gorm:"not null"`
	Line   int    `gorm:"not null"`
}

type ScriptRequest struct {
	Src     string `json:"src"`
	Content string `json:"content"`
}

type FindingResponse struct {
	ID     uint `json:"id"`
	ScanID uint `json:"scan_id"`

	RegexName        string `json:"regex_name"`
	RegexDescription string `json:"regex_description"`

	Match  string `json:"match"`
	Source string `json:"source"`
	Line   int    `json:"line"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ConvertFindings(findings []FindingModel) []FindingResponse {
	var findingResponses []FindingResponse

	for _, finding := range findings {
		findingResponses = append(findingResponses, FindingResponse{
			ID:     finding.ID,
			ScanID: finding.ScanID,
			Line:   finding.Line,
			Match:  finding.Match,
			Source: finding.Source,

			RegexName:        finding.RegexName,
			RegexDescription: finding.RegexDescription,

			CreatedAt: finding.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: finding.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return findingResponses
}

func FindFinding(findings []FindingModel, finding FindingModel) bool {
	for _, f := range findings {
		if f.Line == finding.Line && f.Match == finding.Match && f.Source == finding.Source {
			return true
		}
	}

	return false
}

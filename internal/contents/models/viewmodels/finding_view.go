package viewmodels

import (
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
)

type Finding struct {
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

func ConvertFindings(findings []entities.FindingModel) []Finding {
	var findingResponses []Finding

	for _, finding := range findings {

		findingResponses = append(findingResponses, Finding{
			ID:     finding.Id,
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

func FindFinding(findings []entities.FindingModel, finding entities.FindingModel) bool {
	for _, f := range findings {
		if f.Line == finding.Line && f.Match == finding.Match && f.Source == finding.Source {
			return true
		}
	}

	return false
}

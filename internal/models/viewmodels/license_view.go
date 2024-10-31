package viewmodels

import "github.com/codevault-llc/humblebrag-api/internal/models/entities"

func ConvertLicense(license entities.LicenseModel) License {
	return License{
		ID:        license.ID,
		License:   license.License,
		CreatedAt: license.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: license.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

type License struct {
	ID      uint   `json:"id"`
	License string `json:"license"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

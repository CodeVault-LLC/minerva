package models

import "gorm.io/gorm"

type LicenseModel struct {
	gorm.Model

	License string `gorm:"not null"`

	Scans []ScanModel `gorm:"foreignKey:LicenseID"`
}

type LicenseResponse struct {
	ID uint `json:"id"`

	License string `json:"license"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ConvertLicense(user LicenseModel) LicenseResponse {
	return LicenseResponse{
		ID:        user.ID,
		CreatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

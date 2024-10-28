package models

import (
	"gorm.io/gorm"
)

type CertificateResultModel struct {
	gorm.Model

	CertificateId uint
	Certificate   *CertificateModel `gorm:"foreignKey:CertificateId"`

	Expired         bool `gorm:"not null"`
	Trusted         bool `gorm:"not null"`
	Weak            bool `gorm:"not null"`
	MaliciousIssuer bool `gorm:"not null"`
	Revoked         bool `gorm:"not null"`
	Domain          bool `gorm:"not null"`
}

type CertificateResultModelResponse struct {
	ID uint `json:"id"`

	Expired         bool `json:"expired"`
	Trusted         bool `json:"trusted"`
	Weak            bool `json:"weak"`
	MaliciousIssuer bool `json:"malicious_issuer"`
	Revoked         bool `json:"revoked"`
	Domain          bool `json:"domain"`
}

func ConvertCertificateResult(cert CertificateResultModel) CertificateResultModelResponse {
	return CertificateResultModelResponse{
		ID:              cert.ID,
		Expired:         cert.Expired,
		Trusted:         cert.Trusted,
		Weak:            cert.Weak,
		MaliciousIssuer: cert.MaliciousIssuer,
		Revoked:         cert.Revoked,
		Domain:          cert.Domain,
	}
}

type CertificateResult struct {
	Expired         bool
	Trusted         bool
	Weak            bool
	MaliciousIssuer bool
	Revoked         bool
	Domain          bool
}

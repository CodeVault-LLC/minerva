package viewmodels

import (
	"encoding/json"
	"time"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

// Local type embedding entities.CertificateModel
type CertificateModel struct {
	entities.CertificateModel
}

// Custom methods to handle PublicKey marshaling/unmarshaling
func (c *CertificateModel) BeforeSave(tx *gorm.DB) (err error) {
	if c.PublicKey != "" {
		encodedPublicKey, err := json.Marshal(c.PublicKey)
		if err != nil {
			return err
		}
		c.PublicKey = string(encodedPublicKey)
	}
	return nil
}

func (c *CertificateModel) AfterFind(tx *gorm.DB) (err error) {
	if len(c.PublicKey) > 0 {
		err = json.Unmarshal([]byte(c.PublicKey), &c.PublicKey)
		if err != nil {
			return err
		}
	}
	return nil
}

type Certificate struct {
	ID uint `json:"id"`

	Subject string `json:"subject"`
	Issuer  string `json:"issuer"`

	NotBefore time.Time `json:"not_before"`
	NotAfter  time.Time `json:"not_after"`

	SignatureAlgorithm string `json:"signature_algorithm"`
	PublicKeyAlgorithm string `json:"public_key_algorithm"`
}

func ConvertCertificate(certificate entities.CertificateModel) Certificate {
	return Certificate{
		ID:                 certificate.ID,
		Issuer:             certificate.Issuer,
		Subject:            certificate.Subject,
		NotBefore:          certificate.NotBefore,
		NotAfter:           certificate.NotAfter,
		SignatureAlgorithm: certificate.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: certificate.PublicKeyAlgorithm.String(),
	}
}

func ConvertCertificates(certificates []entities.CertificateModel) []Certificate {
	var certificateResponses []Certificate

	for _, certificate := range certificates {
		certificateResponses = append(certificateResponses, Certificate{
			ID:                 certificate.ID,
			Issuer:             certificate.Issuer,
			Subject:            certificate.Subject,
			NotBefore:          certificate.NotBefore,
			NotAfter:           certificate.NotAfter,
			SignatureAlgorithm: certificate.SignatureAlgorithm.String(),
			PublicKeyAlgorithm: certificate.PublicKeyAlgorithm.String(),
		})
	}

	if len(certificateResponses) == 0 {
		return []Certificate{}
	}

	return certificateResponses
}

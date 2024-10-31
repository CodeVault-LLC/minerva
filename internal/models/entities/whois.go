package entities

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type WhoisModel struct {
	gorm.Model

	NetworkId uint
	Network   *NetworkModel `gorm:"foreignKey:NetworkId"`

	DomainName  string         `gorm:"not null"`
	Registrar   string         `gorm:"not null"`
	Email       string         `gorm:"not null"`
	Phone       string         `gorm:"not null"`
	Updated     string         `gorm:"not null"`
	Created     string         `gorm:"not null"`
	Expires     string         `gorm:"not null"`
	Status      string         `gorm:"not null"`
	NameServers pq.StringArray `gorm:"type:text[]"`

	RegistrantName       string `gorm:"not null"`
	RegistrantEmail      string `gorm:"not null"`
	RegistrantPhone      string `gorm:"not null"`
	RegistrantOrg        string `gorm:"not null"`
	RegistrantCity       string `gorm:"not null"`
	RegistrantCountry    string `gorm:"not null"`
	RegistrantPostalCode string `gorm:"not null"`

	AdminName       string `gorm:"not null"`
	AdminEmail      string `gorm:"not null"`
	AdminPhone      string `gorm:"not null"`
	AdminOrg        string `gorm:"not null"`
	AdminCity       string `gorm:"not null"`
	AdminCountry    string `gorm:"not null"`
	AdminPostalCode string `gorm:"not null"`
}


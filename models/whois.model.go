package models

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

type WhoisResponse struct {
	ID uint `json:"id"`

	DomainName  string   `json:"domain_name"`
	Registrar   string   `json:"registrar"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Updated     string   `json:"updated"`
	Created     string   `json:"created"`
	Expires     string   `json:"expires"`
	Status      string   `json:"status"`
	NameServers []string `json:"name_servers"`

	RegistrantName       string `json:"registrant_name"`
	RegistrantEmail      string `json:"registrant_email"`
	RegistrantPhone      string `json:"registrant_phone"`
	RegistrantOrg        string `json:"registrant_org"`
	RegistrantCity       string `json:"registrant_city"`
	RegistrantState      string `json:"registrant_state"`
	RegistrantCountry    string `json:"registrant_country"`
	RegistrantPostalCode string `json:"registrant_postal_code"`

	AdminName       string `json:"admin_name"`
	AdminEmail      string `json:"admin_email"`
	AdminPhone      string `json:"admin_phone"`
	AdminOrg        string `json:"admin_org"`
	AdminCity       string `json:"admin_city"`
	AdminState      string `json:"admin_state"`
	AdminCountry    string `json:"admin_country"`
	AdminPostalCode string `json:"admin_postal_code"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ConvertWhois(whois WhoisModel) WhoisResponse {
	return WhoisResponse{
		ID: whois.ID,

		DomainName:  whois.DomainName,
		Registrar:   whois.Registrar,
		Email:       whois.Email,
		Phone:       whois.Phone,
		Updated:     whois.Updated,
		Created:     whois.Created,
		Expires:     whois.Expires,
		Status:      whois.Status,
		NameServers: whois.NameServers,

		RegistrantName:       whois.RegistrantName,
		RegistrantEmail:      whois.RegistrantEmail,
		RegistrantPhone:      whois.RegistrantPhone,
		RegistrantOrg:        whois.RegistrantOrg,
		RegistrantCity:       whois.RegistrantCity,
		RegistrantState:      whois.RegistrantCountry,
		RegistrantCountry:    whois.RegistrantCountry,
		RegistrantPostalCode: whois.RegistrantPostalCode,

		AdminName:       whois.AdminName,
		AdminEmail:      whois.AdminEmail,
		AdminPhone:      whois.AdminPhone,
		AdminOrg:        whois.AdminOrg,
		AdminCity:       whois.AdminCity,
		AdminState:      whois.AdminCountry,
		AdminCountry:    whois.AdminCountry,
		AdminPostalCode: whois.AdminPostalCode,

		CreatedAt: whois.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: whois.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

package viewmodels

import "github.com/codevault-llc/humblebrag-api/internal/models/entities"

type Whois struct {
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

func ConvertWhois(whois entities.WhoisModel) Whois {
	return Whois{
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

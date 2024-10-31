package viewmodels

import "github.com/codevault-llc/humblebrag-api/internal/models/entities"

type Metadata struct {
	ID uint `json:"id"`

	Robots  string `json:"robots"`
	Readme  string `json:"readme"`
	License string `json:"license"`

	CMS            string   `json:"cms"`
	ServerSoftware string   `json:"server_software"`
	Frameworks     []string `json:"frameworks"`
	ServerLanguage string   `json:"server_language"`
}

func ConvertMetadata(metadata entities.MetadataModel) Metadata {
	return Metadata{
		ID: metadata.ID,

		Robots:  metadata.Robots,
		Readme:  metadata.Readme,
		License: metadata.License,

		CMS:            metadata.CMS,
		ServerSoftware: metadata.ServerSoftware,
		Frameworks:     metadata.Frameworks,
	}
}

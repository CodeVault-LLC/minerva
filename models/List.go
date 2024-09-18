package models

import (
	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/types"
	"gorm.io/gorm"
)

type List struct {
	gorm.Model

	ScanID uint
	Scan   Scan

	ListID string // match towards the listID in the config
}

type ListResponse struct {
	ID uint `json:"id"`

	Description string   `json:"description"`
	ListID      string   `json:"list_id"`
	Categories  []string `json:"categories"`
	URL         string   `json:"url"`
}

func ConvertList(list List) ListResponse {
	var configList *types.List
	for _, l := range config.ConfigLists {
		if l.ListID == list.ListID {
			configList = l
			break
		}
	}

	return ListResponse{
		ID:          list.ID,
		Description: configList.Description,
		ListID:      list.ListID,
		Categories:  configList.Categories,
		URL:         configList.URL,
	}
}

func ConvertLists(lists []List) []ListResponse {
	var listResponses []ListResponse

	for _, list := range lists {
		listResponses = append(listResponses, ConvertList(list))
	}

	return listResponses
}

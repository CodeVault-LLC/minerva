package models

import (
	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/types"
	"gorm.io/gorm"
)

type FilterModel struct {
	gorm.Model

	ScanID uint
	Scan   ScanModel

	FilterID string // match towards the listID in the config
}

type FilterResponse struct {
	ID uint `json:"id"`

	Description string   `json:"description"`
	FilterID    string   `json:"filter_id"`
	Categories  []string `json:"categories"`
	URL         string   `json:"url"`
}

func ConvertFilter(list FilterModel) FilterResponse {
	var configList *types.Filter
	for _, l := range config.ConfigLists {
		if l.FilterID == list.FilterID {
			configList = l
			break
		}
	}

	return FilterResponse{
		ID:          list.ID,
		Description: configList.Description,
		FilterID:    list.FilterID,
		Categories:  configList.Categories,
		URL:         configList.URL,
	}
}

func ConvertFilters(lists []FilterModel) []FilterResponse {
	var filterResponses []FilterResponse

	for _, list := range lists {
		filterResponses = append(filterResponses, ConvertFilter(list))
	}

	return filterResponses
}

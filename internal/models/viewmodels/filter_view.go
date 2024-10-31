package viewmodels

import (
	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
)

type Filter struct {
	ID uint `json:"id"`

	Description string   `json:"description"`
	FilterID    string   `json:"filter_id"`
	Categories  []string `json:"categories"`
	URL         string   `json:"url"`
}

func ConvertFilter(list entities.FilterModel) Filter {
	var configList *types.Filter
	for _, l := range config.ConfigLists {
		if l.FilterID == list.FilterID {
			configList = l
			break
		}
	}

	return Filter{
		ID:          list.ID,
		Description: configList.Description,
		FilterID:    list.FilterID,
		Categories:  configList.Categories,
		URL:         configList.URL,
	}
}

func ConvertFilters(lists []entities.FilterModel) []Filter {
	var filterResponses []Filter

	for _, list := range lists {
		filterResponses = append(filterResponses, ConvertFilter(list))
	}

	return filterResponses
}

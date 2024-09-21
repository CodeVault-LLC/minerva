package list

import (
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/parsers"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
)

func ListModule(scanId uint, url string) {
	foundLists := updater.CompareValues(utils.ConvertURLToDomain(url), parsers.Domain)

	for _, list := range foundLists {
		listModel := models.ListModel{
			ScanID: scanId,
			ListID: list.ListID,
		}

		_, err := service.CreateList(listModel)
		if err != nil {
			logger.Log.Error("Failed to save list: %v", err)
		}
	}
}

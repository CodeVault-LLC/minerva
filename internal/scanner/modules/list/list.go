package list

import (
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/parsers"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"go.uber.org/zap"
)

func ListModule(scanId uint, url string) {
	foundLists := updater.CompareValues(utils.ConvertURLToDomain(url), parsers.Domain)

	for _, list := range foundLists {
		filterModel := models.FilterModel{
			ScanID:   scanId,
			FilterID: list.FilterID,
		}

		_, err := service.CreateFilter(filterModel)
		if err != nil {
			logger.Log.Error("Failed to save list: %v", zap.Error(err))
		}
	}
}

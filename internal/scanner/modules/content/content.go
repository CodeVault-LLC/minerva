package content

import (
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
)

func ContentModule(scanId uint, scripts []models.ScriptRequest) {
	for _, script := range scripts {
		content := models.ContentModel{
			ScanID: scanId,

			Name:    script.Src,
			Content: script.Content,
		}

		_, err := service.CreateContent(content)
		if err != nil {
			logger.Log.Error("Failed to save content: %v", err)
		}
	}
}

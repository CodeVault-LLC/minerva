package models

import "gorm.io/gorm"

type ContentModel struct {
	gorm.Model

	ScanID uint
	Scan   ScanModel

	Name string `gorm:"not null"`

	Content string `gorm:"not null"`
}

type ContentResponse struct {
	ID      uint   `json:"id"`
	ScanID  uint   `json:"scan_id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func ConvertContents(content []ContentModel) []ContentResponse {
	var contentResponses []ContentResponse

	for _, c := range content {
		contentResponses = append(contentResponses, ContentResponse{
			ID:      c.ID,
			ScanID:  c.ScanID,
			Name:    c.Name,
			Content: c.Content,
		})
	}

	return contentResponses
}

func ConvertContent(content ContentModel) ContentResponse {
	return ContentResponse{
		ID:      content.ID,
		ScanID:  content.ScanID,
		Name:    content.Name,
		Content: content.Content,
	}
}

package viewmodels

import (
	"github.com/codevault-llc/minerva/internal/contents/models/entities"
)

// Contents represents the data returned in API responses.
type Contents struct {
	ID          string   `json:"id"`
	FileSize    int64    `json:"file_size"`
	FileType    string   `json:"file_type"`
	StorageType string   `json:"storage_type"`
	AccessCount int64    `json:"access_count"`
	Tags        []string `json:"tags"`
	ObjectKey   string   `json:"object_key"`
}

type Content struct {
	ID          string `json:"id"`
	FileSize    int64  `json:"file_size"`
	FileType    string `json:"file_type"`
	StorageType string `json:"storage_type"`
	AccessCount int64  `json:"access_count"`
	Body        string `json:"body"`
}

// ConvertContents converts a list of ContentModel to ContentResponse.
func ConvertContents(contents []entities.ContentModel, storageMap map[string]entities.ContentStorageModel) []Contents {
	var contentResponses []Contents

	for _, c := range contents {
		contentResponses = append(contentResponses, ConvertContent(c, storageMap[c.Id]))
	}

	return contentResponses
}

// ConvertContent converts a ContentModel to ContentResponse.
func ConvertContent(content entities.ContentModel, storage entities.ContentStorageModel) Contents {
	return Contents{
		ID:          content.Id,
		FileSize:    content.FileSize,
		FileType:    content.FileType,
		StorageType: content.StorageType,
		ObjectKey:   storage.ObjectKey,
	}
}

func ConvertSingleContent(content entities.ContentModel) Contents {
	return Contents{
		ID:          content.Id,
		FileSize:    content.FileSize,
		FileType:    content.FileType,
		StorageType: content.StorageType,
	}
}

package viewmodels

import (
	"time"

	"github.com/codevault-llc/humblebrag-api/internal/contents/models/entities"
)

// Contents represents the data returned in API responses.
type Contents struct {
	ID           uint      `json:"id"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	StorageType  string    `json:"storage_type"`
	LastAccessed time.Time `json:"last_accessed"`
	AccessCount  int64     `json:"access_count"`
	Tags         []string  `json:"tags"`
	ObjectKey    string    `json:"object_key"`
}

type Content struct {
	ID           uint      `json:"id"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	StorageType  string    `json:"storage_type"`
	LastAccessed time.Time `json:"last_accessed"`
	AccessCount  int64     `json:"access_count"`
	Body         string    `json:"body"`
}

// ConvertContents converts a list of ContentModel to ContentResponse.
func ConvertContents(contents []entities.ContentModel, tagsMap map[uint][]string, storageMap map[uint]entities.ContentStorageModel) []Contents {
	var contentResponses []Contents

	for _, c := range contents {
		contentResponses = append(contentResponses, ConvertContent(c, tagsMap[c.ID], storageMap[c.ID]))
	}

	return contentResponses
}

// ConvertContent converts a ContentModel to ContentResponse.
func ConvertContent(content entities.ContentModel, tags []string, storage entities.ContentStorageModel) Contents {
	return Contents{
		ID:           content.ID,
		FileSize:     content.FileSize,
		FileType:     content.FileType,
		StorageType:  content.StorageType,
		LastAccessed: content.LastAccessedAt,
		AccessCount:  content.AccessCount,
		Tags:         tags,
		ObjectKey:    storage.ObjectKey,
	}
}

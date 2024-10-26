package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentModel represents the core information about a piece of content.
type ContentModel struct {
	gorm.Model

	ID uint `gorm:"primaryKey"` // Unique identifier for the content

	HashedBody     string    `gorm:"type:varchar(255);not null;unique"` // Hash of the content body for deduplication
	SourceID       uint      `gorm:"not null"`                          // Reference to the source (e.g., a website or client)
	FileSize       int64     `gorm:"not null"`                          // Size of the file in bytes
	FileType       string    `gorm:"type:varchar(100);not null"`        // MIME type of the content (e.g., "text/plain", "image/jpeg")
	StorageType    string    `gorm:"type:varchar(50);not null"`         // "hot" or "cold" indicating storage tier
	LastAccessedAt time.Time // Timestamp for the last access, useful for determining storage transitions
	AccessCount    int64     `gorm:"not null;default:0"` // Count of how many times the content has been accessed
}

// ContentStorageModel holds information about where the content is stored.
type ContentStorageModel struct {
	gorm.Model

	ContentID       uint   `gorm:"not null;index"`             // Reference to ContentModel
	BucketName      string `gorm:"type:varchar(255);not null"` // S3 bucket name
	ObjectKey       string `gorm:"type:varchar(255);not null"` // Key of the object in the bucket
	Location        string `gorm:"type:varchar(255);not null"` // Region or specific location if needed (e.g., "us-east-1")
	StorageEndpoint string `gorm:"type:varchar(255);not null"` // Endpoint URL for accessing the storage (e.g., S3 or MinIO endpoint)
	Encryption      string `gorm:"type:varchar(50);not null"`  // Encryption method used (e.g., "AES256", "none")
}

// ContentTagsModel stores tags for categorizing and searching content.
type ContentTagsModel struct {
	gorm.Model

	ContentID uint   `gorm:"not null;index"`                   // Reference to ContentModel
	Tag       string `gorm:"type:varchar(100);not null;index"` // Tag or label for the content
}

// ContentAccessLogModel logs access events for the content.
type ContentAccessLogModel struct {
	gorm.Model

	ContentID  uint      `gorm:"not null;index"`            // Reference to ContentModel
	AccessedAt time.Time `gorm:"not null"`                  // Timestamp of the access event
	AccessType string    `gorm:"type:varchar(50);not null"` // Type of access (e.g., "read", "download")
	IPAddress  string    `gorm:"type:varchar(45)"`          // IP address of the accessing client
}

// ContentResponse represents the data returned in API responses.
type ContentResponse struct {
	ID           uint      `json:"id"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	StorageType  string    `json:"storage_type"`
	LastAccessed time.Time `json:"last_accessed"`
	AccessCount  int64     `json:"access_count"`
	Tags         []string  `json:"tags"`
	ObjectKey    string    `json:"object_key"`
}

// ConvertContents converts a list of ContentModel to ContentResponse.
func ConvertContents(contents []ContentModel, tagsMap map[uint][]string, storageMap map[uint]ContentStorageModel) []ContentResponse {
	var contentResponses []ContentResponse

	for _, c := range contents {
		contentResponses = append(contentResponses, ConvertContent(c, tagsMap[c.ID], storageMap[c.ID]))
	}

	return contentResponses
}

// ConvertContent converts a ContentModel to ContentResponse.
func ConvertContent(content ContentModel, tags []string, storage ContentStorageModel) ContentResponse {
	return ContentResponse{
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

package entities

import (
	"time"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type ContentModel struct {
	gorm.Model

	ID uint `gorm:"primaryKey"` // Unique identifier for the content

	ScanID uint `gorm:"not null;index"` // Reference to the scan that discovered the content
	Scan   entities.ScanModel

	HashedBody     string    `gorm:"type:varchar(255);not null;unique"` // Hash of the content body for deduplication
	Source         string    `gorm:"type:text;not null"`                // Source URL or origin of the content
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

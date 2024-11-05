package entities

import "time"

type ContentModel struct {
	Id uint `db:"id"`

	ScanId uint `db:"scan_id"`

	HashedBody     string    `db:"hashed_body"`
	Source         string    `db:"source"`
	FileSize       int64     `db:"file_size"`
	FileType       string    `db:"file_type"`
	StorageType    string    `db:"storage_type"`
	LastAccessedAt time.Time `db:"last_accessed_at"`
	AccessCount    int64     `db:"access_count"`

	// Relationships
	Tags     []ContentTagsModel  `db:"-"`
	Storage  ContentStorageModel `db:"-"`
	Access   []ContentAccessLog  `db:"-"`
	Findings []FindingModel      `db:"-"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type ContentTagsModel struct {
	Id uint `db:"id"`

	ContentId uint `db:"content_id"`

	Tag string `db:"tag"`
}

type ContentStorageModel struct {
	Id uint `db:"id"`

	ContentId       uint   `db:"content_id"`
	BucketName      string `db:"bucket_name"`
	ObjectKey       string `db:"object_key"`
	Location        string `db:"location"`
	StorageEndpoint string `db:"storage_endpoint"`
	Encryption      string `db:"encryption"`
}

type ContentAccessLog struct {
	Id uint `db:"id"`

	ContentId uint `db:"content_id"`

	AccessedAt time.Time `db:"accessed_at"`
	AccessType string    `db:"access_type"`
	IpAddress  string    `db:"ip_address"`
}

type FindingModel struct {
	Id uint `db:"id"`

	ScanId uint `db:"scan_id"`

	RegexName        string `db:"regex_name"`
	RegexDescription string `db:"regex_description"`

	Match  string `db:"match"`
	Source string `db:"source"`
	Line   int    `db:"line"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

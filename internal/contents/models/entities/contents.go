package entities

import "time"

type ContentModel struct {
	Id     string `cql:"id"`
	ScanId string `cql:"scan_id"`

	HashedBody  string `cql:"hashed_body"`
	Source      string `cql:"source"`
	FileSize    int64  `cql:"file_size"`
	FileType    string `cql:"file_type"`
	StorageType string `cql:"storage_type"`

	Tags []string `cql:"tags"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

type ContentStorageModel struct {
	Id string `cql:"id"`

	ContentId       string `cql:"content_id"`
	BucketName      string `cql:"bucket_name"`
	ObjectKey       string `cql:"object_key"`
	Location        string `cql:"location"`
	StorageEndpoint string `cql:"storage_endpoint"`
	Encryption      string `cql:"encryption"`
}

type FindingModel struct {
	Id     string `cql:"id"`
	ScanId string `cql:"scan_id"`

	RegexName        string `cql:"regex_name"`
	RegexDescription string `cql:"regex_description"`

	Match  string `cql:"match"`
	Source string `cql:"source"`
	Line   int    `cql:"line"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

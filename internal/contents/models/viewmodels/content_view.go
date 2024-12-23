package viewmodels

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

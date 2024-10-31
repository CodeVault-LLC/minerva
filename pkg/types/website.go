package types

type WebsiteAnalysis struct {
	Url        string        `json:"url"`
	Title      string        `json:"name"`
	StatusCode int           `json:"status_code"`
	Assets     []FileRequest `json:"files"`
	Redirects  []string      `json:"redirects"`
}

type FileRequest struct {
	Src        string `json:"src"`
	HashedBody string `json:"hashed_body"`
	FileSize   uint   `json:"file_size"`
	FileType   string `json:"file_type"`
	Content    string `json:"content"`
}

type Screenshot struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

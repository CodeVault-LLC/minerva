package types

type WebsiteAnalysis struct {
	Url        string        `json:"url"`
	Title      string        `json:"name"`
	StatusCode int           `json:"status_code"`
	Assets     []FileRequest `json:"files"`
	Redirects  []Redirect    `json:"redirects"`
}

type Redirect struct {
	Url        string     `json:"url"`
	StatusCode int        `json:"status_code"`
	Screenshot Screenshot `json:"screenshot"`
}

type FileRequest struct {
	Src        string `json:"src"`
	HashedBody string `json:"hashed_body"`
	FileSize   uint   `json:"file_size"`
	FileType   string `json:"file_type"`
	Content    string `json:"content"`
}

type Screenshot struct {
	Content string `json:"content"`
}

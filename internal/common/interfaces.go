package common

import "github.com/codevault-llc/minerva/internal/core/models/entities"

// ScanModule interface for modules that can execute a scan
type ScanModule interface {
	Execute(job entities.JobModel, website *WebsiteAnalysis) error
	Name() string
	RuntimeLocation() RuntimeLocation
}

// MiniModule interface for smaller sub-modules within a module
type MiniModule interface {
	Run(job entities.JobModel) (interface{}, error)
	Name() string
}

type RuntimeLocation string

// Pre-Scan -> In-Fetch -> Post-Fetch -> Post-Scan
const (
	RuntimeLocationPreScan   RuntimeLocation = "pre-scan"
	RuntimeLocationInFetch   RuntimeLocation = "in-fetch"
	RuntimeLocationPostFetch RuntimeLocation = "post-fetch"
	RuntimeLocationPostScan  RuntimeLocation = "post-scan"
)

type WebsiteAnalysis struct {
	Url       string        `json:"url"`
	Title     string        `json:"name"`
	Assets    []FileRequest `json:"files"`
	Redirects []Redirect    `json:"redirects"`
}

type Redirect struct {
	Url        string     `json:"url"`
	StatusCode int        `json:"status_code"`
	Screenshot Screenshot `json:"screenshot"`
}

type FileRequest struct {
	Src      string `json:"src"`
	FileSize uint   `json:"file_size"`
	FileType string `json:"file_type"`
	Content  string `json:"content"`
	Duration int    `json:"duration"`
	Headers  string `json:"headers"`
	Cookies  string `json:"cookies"`
	Status   int    `json:"status"`
}

type Screenshot struct {
	Content string `json:"content"`
}

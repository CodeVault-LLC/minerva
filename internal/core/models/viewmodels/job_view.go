package viewmodels

import (
	"time"

	"github.com/codevault-llc/humblebrag-api/internal/core/models/entities"
)

// ScanRequest represents the structure of the incoming scan request
type ScanRequest struct {
	URL       string `json:"url"`
	UserAgent string `json:"user_agent,omitempty"`
}

// ConvertJob converts an entity job to a viewmodel job response
func ConvertJob(job entities.JobModel) JobResponse {
	return JobResponse{
		ID:        job.ID,
		ScanID:    job.ScanID,
		Type:      job.Type,
		URL:       job.URL,
		Status:    job.Status,
		CreatedAt: job.CreatedAt,
	}
}

// JobResponse represents the structure for outgoing job data in the API response
type JobResponse struct {
	ID        string             `json:"id"`
	Type      string             `json:"type"`
	ScanID    uint               `json:"scan_id"`
	URL       string             `json:"url"`
	Status    entities.JobStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
}

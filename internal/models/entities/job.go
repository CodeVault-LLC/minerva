package entities

import "time"

// Job represents a task to be processed by the TaskScheduler
type Job struct {
	ID        string
	Type      string
	URL       string
	UserAgent string
	LicenseID int
	Status    JobStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// JobStatus indicates the current status of a Job
type JobStatus int

const (
	Queued JobStatus = iota
	Processing
	Completed
	Failed
)

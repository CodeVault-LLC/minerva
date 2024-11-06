package entities

import (
	"time"
)

// JobModel represents a task to be processed by the TaskScheduler
type JobModel struct {
	ID   string
	Type string

	URL       string
	UserAgent string

	ScanID uint

	Status      JobStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
}

// JobStatus indicates the current status of a Job
type JobStatus int

const (
	Queued JobStatus = iota
	Processing
	Completed
	Failed
)

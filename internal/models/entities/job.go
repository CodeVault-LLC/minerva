package entities

import (
	"time"

	"gorm.io/gorm"
)

// JobModel represents a task to be processed by the TaskScheduler
type JobModel struct {
	gorm.Model
	Type string `gorm:"not null"`

	URL       string `gorm:"not null"`
	UserAgent string `gorm:"not null"`

	LicenseID int  `gorm:"not null"`
	ScanID    uint `gorm:"not null"`

	Status    JobStatus `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// JobStatus indicates the current status of a Job
type JobStatus int

const (
	Queued JobStatus = iota
	Processing
	Completed
	Failed
)

package entities

import "time"

type ScanModel struct {
	Id int `db:"id"`

	Url        string `db:"url"`
	Title      string `db:"title"`
	StatusCode int    `db:"status_code"`

	Status ScanStatus `db:"status"`

	Sha256 string `db:"sha256"`
	Sha1   string `db:"sha1"`
	Md5    string `db:"md5"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type ScanStatus string

const (
	ScanStatusArchived ScanStatus = "archived"
	ScanStatusComplete ScanStatus = "complete"
	ScanStatusFailed   ScanStatus = "failed"
	ScanStatusPending  ScanStatus = "pending"
)

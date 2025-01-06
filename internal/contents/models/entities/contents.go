package entities

import (
	"time"
)

type ContentModel struct {
	Id     string `cql:"id"`
	ScanId string `cql:"scan_id"`

	FileSize int64  `cql:"file_size"` // in bytes
	FileType string `cql:"file_type"` // ex: js, css, html
	Source   string `cql:"source"`    // ex: https://hosting.com/file.js
	Md5      string `cql:"md5"`       // hash of the content
	Sha1     string `cql:"sha1"`      // hash of the content
	Sha256   string `cql:"sha256"`    // hash of the content

	Duration int `cql:"duration"` // in milliseconds

	Tags []string `cql:"tags"`

	Headers string `cql:"headers"`
	Cookies string `cql:"cookies"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

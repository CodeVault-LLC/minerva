package entities

import "time"

type ScanModel struct {
	Id string `cql:"id"`

	Url        string `cql:"url"`
	Title      string `cql:"title"`
	StatusCode int    `cql:"status_code"`

	Sha256 string `cql:"sha256"`
	Sha1   string `cql:"sha1"`
	Md5    string `cql:"md5"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

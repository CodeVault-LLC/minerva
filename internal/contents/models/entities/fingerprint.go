package entities

import "time"

type FingerprintModel struct {
	Id        string `cql:"id"`
	ContentId string `cql:"content_id"`

	FingerprintName string `cql:"fingerprint_name"`

	Match  string `cql:"match"`
	Source string `cql:"source"`
	Line   int    `cql:"line"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

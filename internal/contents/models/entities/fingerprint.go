package entities

import "time"

type FingerprintModel struct {
	Id        string `cql:"id"`
	ContentId string `cql:"content_id"`

	FingerprintId string `cql:"fingerprind_id"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

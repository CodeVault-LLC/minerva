package entities

import "time"

type FindingModel struct {
	Id     string `cql:"id"`
	ScanId string `cql:"scan_id"`

	RegexName        string `cql:"regex_name"`
	RegexDescription string `cql:"regex_description"`

	Match  string `cql:"match"`
	Source string `cql:"source"`
	Line   int    `cql:"line"`

	CreatedAt time.Time `cql:"created_at"`
	UpdatedAt time.Time `cql:"updated_at"`
}

package models

var NetworkSchema = []string{`CREATE TABLE IF NOT EXISTS minerva.networks (
	id UUID PRIMARY KEY,
	scan_id TEXT,
	ip_addresses TEXT,
	ip_ranges TEXT,
	http_headers TEXT,
	dns TEXT,
	whois TEXT,
	certificates TEXT,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);`}

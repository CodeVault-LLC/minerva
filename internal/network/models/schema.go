package models

var NetworkSchema = []string{`CREATE TABLE IF NOT EXISTS networks (
	id UUID DEFAULT generateUUIDv4(),
	scan_id UUID,
	ip_addresses String,
	ip_ranges String,
	http_headers String,
	dns String,
	whois String,
	certificates String,
	created_at DateTime DEFAULT now(),
	updated_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY id;`}

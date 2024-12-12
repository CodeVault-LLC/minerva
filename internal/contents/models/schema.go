package schema

var ContentSchema = []string{
	`CREATE TABLE IF NOT EXISTS minerva.content (
		id UUID PRIMARY KEY,
		scan_id TEXT,

		hashed_body TEXT,
		source TEXT,
		file_size BIGINT,
		file_type TEXT,
		storage_type TEXT,
		tags TEXT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
	);`,
	`CREATE TABLE IF NOT EXISTS minerva.content_storage (
		id UUID PRIMARY KEY,
		content_id TEXT,
		bucket_name TEXT,
		object_key TEXT,
		location TEXT,
		storage_endpoint TEXT,
		encryption TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS minerva.content_access_log (
		id UUID PRIMARY KEY,
		content_id TEXT,
		accessed_at TIMESTAMP,
		access_type TEXT,
		ip_address TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS minerva.finding (
		id UUID PRIMARY KEY,
		scan_id TEXT,

		regex_name TEXT,
		regex_description TEXT,
		match TEXT,
		source TEXT,
		line INT,

		created_at TIMESTAMP,
		updated_at TIMESTAMP,
	);`,
	`CREATE INDEX IF NOT EXISTS idx_hashed_body ON minerva.content (hashed_body);`,
	`CREATE INDEX IF NOT EXISTS idx_content_id ON minerva.content_storage (content_id);`,
	`CREATE INDEX IF NOT EXISTS idx_scan_id ON minerva.content (scan_id);`,
}

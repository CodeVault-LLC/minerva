package schema

var ContentSchema = []string{
	`CREATE TABLE IF NOT EXISTS minerva.content (
		id TEXT PRIMARY KEY,

		scan_id INT,

		hashed_body TEXT,
		source TEXT,
		file_size BIGINT,
		file_type TEXT,
		storage_type TEXT,
		last_accessed_at TIMESTAMP,
		access_count BIGINT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP
	);`,
	`CREATE TABLE IF NOT EXISTS minerva.content_storage (
		id TEXT PRIMARY KEY,
		content_id INT,
		bucket_name TEXT,
		object_key TEXT,
		location TEXT,
		storage_endpoint TEXT,
		encryption TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS minerva.content_tags (
		id TEXT PRIMARY KEY,
		content_id INT,
		tag TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS minerva.content_access_log (
		id TEXT PRIMARY KEY,
		content_id INT,
		accessed_at TIMESTAMP,
		access_type TEXT,
		ip_address TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS minerva.finding (
		id TEXT PRIMARY KEY,
		scan_id INT,

		regex_name TEXT,
		regex_description TEXT,
		match TEXT,
		source TEXT,
		line INT,

		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP
	);`,
}

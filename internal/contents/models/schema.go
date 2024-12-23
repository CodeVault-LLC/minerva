package schema

var ContentSchema = []string{
	`CREATE TABLE IF NOT EXISTS content (
		id UUID,
		scan_id UUID,
		file_size Int32,
		file_type String,
		source String,
		md5 String,
		sha1 String,
		sha256 String,
		duration Int32,
		tags Array(String),
		created_at DateTime DEFAULT now(),
		updated_at DateTime DEFAULT now()
	) ENGINE = MergeTree()
	ORDER BY id;`,
	`CREATE TABLE IF NOT EXISTS finding (
		id UUID DEFAULT generateUUIDv4(),
		scan_id UUID,

		regex_name String,
		regex_description String,
		match String,
		source String,
		line Int32,

		created_at DateTime DEFAULT now(),
		updated_at DateTime DEFAULT now()
	) ENGINE = MergeTree()
	ORDER BY id;`,
	`CREATE TABLE IF NOT EXISTS fingerprint (
		id UUID DEFAULT generateUUIDv4(),
		content_id UUID,
		fingerprint_name String,
		match String,
		source String,
		line Int32,
		created_at DateTime DEFAULT now(),
		updated_at DateTime DEFAULT now()
	) ENGINE = MergeTree()
	ORDER BY id;`,
}

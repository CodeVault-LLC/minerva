package schema

var ContentSchema = `
CREATE TABLE IF NOT EXISTS content (
	id SERIAL PRIMARY KEY,

	-- Scan relationship
	scan_id INTEGER NOT NULL,
	FOREIGN KEY (scan_id) REFERENCES scans(id) ON DELETE CASCADE,

	-- Content details
	hashed_body VARCHAR(255) NOT NULL UNIQUE, -- Unique hash of the content body
	source TEXT NOT NULL,                     -- Source URL or origin
	file_size BIGINT NOT NULL,                -- File size in bytes
	file_type VARCHAR(100) NOT NULL,          -- MIME type of the content
	storage_type VARCHAR(50) NOT NULL,        -- Storage tier, e.g., "hot" or "cold"
	last_accessed_at TIMESTAMP,               -- Timestamp for the last access
	access_count BIGINT NOT NULL DEFAULT 0,   -- Access count, default to 0

	-- Timestamps for created_at, updated_at, deleted_at
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS content_storage (
	id SERIAL PRIMARY KEY,

	-- Content relationship
	content_id INTEGER NOT NULL,
	FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,

	-- Storage details
	bucket_name VARCHAR(255) NOT NULL,       -- S3 bucket name
	object_key VARCHAR(255) NOT NULL,        -- Key of the object in the bucket
	location VARCHAR(255) NOT NULL,          -- Region or specific location (e.g., "us-east-1")
	storage_endpoint VARCHAR(255) NOT NULL,  -- Endpoint URL for accessing the storage
	encryption VARCHAR(50) NOT NULL          -- Encryption method (e.g., "AES256")
);

CREATE TABLE IF NOT EXISTS content_tags (
	id SERIAL PRIMARY KEY,

	-- Content relationship
	content_id INTEGER NOT NULL,
	FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,

	-- Tag information
	tag VARCHAR(100) NOT NULL,

	-- Unique constraint to prevent duplicate tags for the same content
	UNIQUE (content_id, tag)
);

CREATE TABLE IF NOT EXISTS content_access_log (
	id SERIAL PRIMARY KEY,

	-- Content relationship
	content_id INTEGER NOT NULL,
	FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,

	-- Access log details
	accessed_at TIMESTAMP NOT NULL,         -- Timestamp of access
	access_type VARCHAR(50) NOT NULL,       -- Type of access (e.g., "read", "download")
	ip_address VARCHAR(45)                  -- IP address of the accessing client
);

CREATE TABLE IF NOT EXISTS finding (
	id SERIAL PRIMARY KEY,

	-- Scan relationship
	scan_id INTEGER NOT NULL,
	FOREIGN KEY (scan_id) REFERENCES scans(id) ON DELETE CASCADE,

	-- Finding details
	regex_name VARCHAR(100) NOT NULL,       -- Name of the regex pattern
	regex_description TEXT NOT NULL,         -- Description of the regex pattern
	match TEXT NOT NULL,                     -- Matched content
	source TEXT NOT NULL,                    -- Source URL or origin
	line INTEGER NOT NULL,                   -- Line number of the match

	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP
);

-- Indexes for optimization
CREATE INDEX IF NOT EXISTS idx_content_scan_id ON content(scan_id);
CREATE INDEX IF NOT EXISTS idx_content_storage_content_id ON content_storage(content_id);
CREATE INDEX IF NOT EXISTS idx_content_tags_content_id ON content_tags(content_id);
CREATE INDEX IF NOT EXISTS idx_content_tags_tag ON content_tags(tag);
CREATE INDEX IF NOT EXISTS idx_content_access_log_content_id ON content_access_log(content_id);
CREATE INDEX IF NOT EXISTS idx_finding_scan_id ON finding(scan_id);
`

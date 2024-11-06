package models

var CoreSchema = `
CREATE TABLE IF NOT EXISTS scans (
	id SERIAL PRIMARY KEY,

	url TEXT NOT NULL,                     -- URL to scan
	title VARCHAR(255),                    -- Title of the page
	status_code INTEGER NOT NULL,          -- HTTP status code

	status VARCHAR(50) NOT NULL DEFAULT 'complete', -- Scan status

	sha256 VARCHAR(64) NOT NULL,            -- SHA256 hash
	sha1 VARCHAR(40) NOT NULL,              -- SHA1 hash
	md5 VARCHAR(32) NOT NULL,               -- MD5 hash

	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP
);
`

package models

var CoreSchema = []string{`CREATE TABLE IF NOT EXISTS minerva.scans (
    id TEXT PRIMARY KEY,

    url TEXT,
    title TEXT,
    status_code INT,

    status TEXT,

    sha256 TEXT,
    sha1 TEXT,
    md5 TEXT,

    created_at TIMESTAMP,
    updated_at TIMESTAMP
);`}

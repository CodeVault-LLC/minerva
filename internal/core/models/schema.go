package models

var CoreSchema = []string{`CREATE TABLE IF NOT EXISTS scans (
    id UUID DEFAULT generateUUIDv4(),

    url String,
    title String,
    status_code UInt16,

    sha256 String,
    sha1 String,
    md5 String,

    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY id;`}

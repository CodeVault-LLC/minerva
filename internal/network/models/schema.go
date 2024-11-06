package models

var NetworkSchema = `
CREATE TABLE IF NOT EXISTS networks (
	id SERIAL PRIMARY KEY,

	-- Scan relationship
	scan_id INTEGER NOT NULL,
	FOREIGN KEY (scan_id) REFERENCES scans(id),

	-- Detail fields
	ip_addresses TEXT[] NOT NULL, -- PostgreSQL array for IP addresses
	ip_ranges TEXT[] NOT NULL,    -- PostgreSQL array for IP ranges

	-- HTTP fields
	http_headers TEXT[] NOT NULL,  -- PostgreSQL array for HTTP headers

	-- Timestamps
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS whois (
	id SERIAL PRIMARY KEY,

	-- Network relationship
	network_id INTEGER NOT NULL,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE,

	-- Whois details
	domain_name TEXT,
	registrar TEXT,
	email TEXT,
	phone TEXT,
	updated TEXT,
	created TEXT,
	expires TEXT,
	status TEXT,
	name_servers TEXT[],  -- PostgreSQL array for name servers

	-- Registrant information
	registrant_name TEXT,
	registrant_email TEXT,
	registrant_phone TEXT,
	registrant_org TEXT,
	registrant_city TEXT,
	registrant_country TEXT,
	registrant_postal_code TEXT,

	-- Admin information
	admin_name TEXT,
	admin_email TEXT,
	admin_phone TEXT,
	admin_org TEXT,
	admin_city TEXT,
	admin_country TEXT,
	admin_postal_code TEXT
);

CREATE TABLE IF NOT EXISTS dns (
	id SERIAL PRIMARY KEY,

	-- Network relationship
	network_id INTEGER NOT NULL,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE,

	-- DNS records
	cname TEXT[],       -- PostgreSQL array for CNAME records
	a_records TEXT[],   -- PostgreSQL array for A records
	aaaa_records TEXT[],-- PostgreSQL array for AAAA records
	mx_records TEXT[],  -- PostgreSQL array for MX records
	ns_records TEXT[],  -- PostgreSQL array for NS records
	txt_records TEXT[], -- PostgreSQL array for TXT records
	ptr_record TEXT,    -- PTR record
	dnssec BOOLEAN      -- DNSSEC enabled/disabled
);

CREATE TABLE IF NOT EXISTS certificates (
	id SERIAL PRIMARY KEY,

	-- Network relationship
	network_id INTEGER,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE,

	-- Certificate details
	subject TEXT,
	issuer TEXT,

	not_before TIMESTAMP,
	not_after TIMESTAMP,

	signature_algorithm TEXT,
	signature BYTEA,  -- Binary signature

	public_key_algorithm TEXT,
	public_key TEXT,  -- JSON string representation of public key

	serial_number TEXT,
	version INTEGER,
	key_usage TEXT,  -- Text representation of key usage

	basic_constraints_valid BOOLEAN,
	is_ca BOOLEAN,

	-- Certificate subject alternative names (SANs)
	dns_names TEXT[],        -- PostgreSQL array for DNS names
	email_addresses TEXT[],  -- PostgreSQL array for email addresses
	ip_addresses TEXT[],     -- PostgreSQL array for IP addresses
	uris TEXT[],             -- PostgreSQL array for URIs

	-- Name constraints
	permitted_dns_domains_critical BOOLEAN,
	permitted_dns_domains TEXT[], -- PostgreSQL array
	excluded_dns_domains TEXT[],  -- PostgreSQL array
	permitted_ip_ranges TEXT[],   -- PostgreSQL array
	excluded_ip_ranges TEXT[],    -- PostgreSQL array
	permitted_email_addresses TEXT[], -- PostgreSQL array
	excluded_email_addresses TEXT[],  -- PostgreSQL array
	permitted_uri_domains TEXT[],     -- PostgreSQL array
	excluded_uri_domains TEXT[]       -- PostgreSQL array
);

-- Indexes for optimization
CREATE INDEX IF NOT EXISTS idx_networks_scan_id ON networks(scan_id);
CREATE INDEX IF NOT EXISTS idx_whois_network_id ON whois(network_id);
CREATE INDEX IF NOT EXISTS idx_dns_network_id ON dns(network_id);
CREATE INDEX IF NOT EXISTS idx_certificates_network_id ON certificates(network_id);
`

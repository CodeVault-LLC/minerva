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
	http_headers TEXT[] NOT NULL  -- PostgreSQL array for HTTP headers
);

CREATE TABLE IF NOT EXISTS whois (
	id SERIAL PRIMARY KEY,

	-- Network relationship
	network_id INTEGER NOT NULL,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE,

	-- Whois details
	domain_name TEXT NOT NULL,
	registrar TEXT NOT NULL,
	email TEXT NOT NULL,
	phone TEXT NOT NULL,
	updated TEXT NOT NULL,
	created TEXT NOT NULL,
	expires TEXT NOT NULL,
	status TEXT NOT NULL,
	name_servers TEXT[] NOT NULL,  -- PostgreSQL array for name servers

	-- Registrant information
	registrant_name TEXT NOT NULL,
	registrant_email TEXT NOT NULL,
	registrant_phone TEXT NOT NULL,
	registrant_org TEXT NOT NULL,
	registrant_city TEXT NOT NULL,
	registrant_country TEXT NOT NULL,
	registrant_postal_code TEXT NOT NULL,

	-- Admin information
	admin_name TEXT NOT NULL,
	admin_email TEXT NOT NULL,
	admin_phone TEXT NOT NULL,
	admin_org TEXT NOT NULL,
	admin_city TEXT NOT NULL,
	admin_country TEXT NOT NULL,
	admin_postal_code TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS dns (
	id SERIAL PRIMARY KEY,

	-- Network relationship
	network_id INTEGER NOT NULL,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE,

	-- DNS records
	cname TEXT[] NOT NULL,       -- PostgreSQL array for CNAME records
	a_records TEXT[] NOT NULL,   -- PostgreSQL array for A records
	aaaa_records TEXT[] NOT NULL,-- PostgreSQL array for AAAA records
	mx_records TEXT[] NOT NULL,  -- PostgreSQL array for MX records
	ns_records TEXT[] NOT NULL,  -- PostgreSQL array for NS records
	txt_records TEXT[] NOT NULL, -- PostgreSQL array for TXT records
	ptr_record TEXT NOT NULL,    -- PTR record
	dnssec BOOLEAN NOT NULL      -- DNSSEC enabled/disabled
);

CREATE TABLE IF NOT EXISTS certificates (
	id SERIAL PRIMARY KEY,

	-- Network relationship
	network_id INTEGER NOT NULL,
	FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE,

	-- Certificate details
	subject TEXT NOT NULL,
	issuer TEXT NOT NULL,

	not_before TIMESTAMP NOT NULL,
	not_after TIMESTAMP NOT NULL,

	signature_algorithm TEXT NOT NULL,
	signature BYTEA NOT NULL,  -- Binary signature

	public_key_algorithm TEXT NOT NULL,
	public_key TEXT NOT NULL,  -- JSON string representation of public key

	serial_number TEXT NOT NULL,
	version INTEGER NOT NULL,
	key_usage TEXT NOT NULL,  -- Text representation of key usage

	basic_constraints_valid BOOLEAN NOT NULL,
	is_ca BOOLEAN NOT NULL,

	-- Certificate subject alternative names (SANs)
	dns_names TEXT[] NOT NULL,        -- PostgreSQL array for DNS names
	email_addresses TEXT[] NOT NULL,  -- PostgreSQL array for email addresses
	ip_addresses TEXT[] NOT NULL,     -- PostgreSQL array for IP addresses
	uris TEXT[] NOT NULL,             -- PostgreSQL array for URIs

	-- Name constraints
	permitted_dns_domains_critical BOOLEAN NOT NULL,
	permitted_dns_domains TEXT[] NOT NULL, -- PostgreSQL array
	excluded_dns_domains TEXT[] NOT NULL,  -- PostgreSQL array
	permitted_ip_ranges TEXT[] NOT NULL,   -- PostgreSQL array
	excluded_ip_ranges TEXT[] NOT NULL,    -- PostgreSQL array
	permitted_email_addresses TEXT[] NOT NULL, -- PostgreSQL array
	excluded_email_addresses TEXT[] NOT NULL,  -- PostgreSQL array
	permitted_uri_domains TEXT[] NOT NULL,     -- PostgreSQL array
	excluded_uri_domains TEXT[] NOT NULL       -- PostgreSQL array
);

-- Indexes for optimization
CREATE INDEX IF NOT EXISTS idx_networks_scan_id ON networks(scan_id);
CREATE INDEX IF NOT EXISTS idx_whois_network_id ON whois(network_id);
CREATE INDEX IF NOT EXISTS idx_dns_network_id ON dns(network_id);
CREATE INDEX IF NOT EXISTS idx_certificates_network_id ON certificates(network_id);
`

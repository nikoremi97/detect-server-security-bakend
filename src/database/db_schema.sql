CREATE TABLE domains (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	domain STRING NULL,
	servers_changed BOOL NULL,
	ssl_grade STRING NULL,
	previous_ssl STRING NULL,
	logo STRING NULL,
	title STRING NULL,
	is_down BOOL NULL,
	updated TIMESTAMPTZ NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	UNIQUE INDEX servers_domain_key (domain ASC),
	FAMILY "primary" (id, domain, servers_changed, ssl_grade, previous_ssl, logo, title, is_down, updated)
);

CREATE TABLE serverdetails (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	address STRING NULL,
	ssl_grade STRING NULL,
	country STRING NULL,
	owner STRING NULL,
	domain_id UUID NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	CONSTRAINT fk_domain_id_ref_domains FOREIGN KEY (domain_id) REFERENCES domains (id),
	INDEX serverdetails_auto_index_fk_domain_id_ref_domains (domain_id ASC),
	FAMILY "primary" (id, address, ssl_grade, country, owner, domain_id)
);
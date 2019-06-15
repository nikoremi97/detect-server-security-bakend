package database

// SQLCreateDomain is the sql query to insert a new domain into domains table
var SQLCreateDomain = `
INSERT INTO domains (
	domain,
	servers_changed,
	ssl_grade,
	previous_ssl,
	logo,
	title,
	is_down,
	updated
)
	VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8
)
	RETURNING id`

// SQLCreateServerDetails is the sql query to insert a new ServerDetail into severdetails table
var SQLCreateServerDetails = `
INSERT INTO serverdetails (
	domain_id,
	address,
	ssl_grade,
	country,
	owner
)
	VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING id`

// SQLGetUpdateTimeDomain is the sql query to get update field from a given domain id
var SQLGetUpdateTimeDomain = `SELECT updated FROM domains WHERE id = $1`

// SQLPreviousSSLGradeDomain is the sql query to get previous_ssl field from a given domain id
var SQLPreviousSSLGradeDomain = `SELECT previous_ssl FROM domains WHERE id = $1`

// SQLCountServerDetails retrives a colum with the number of serverdetails rows that contains a specific domain_id
var SQLCountServerDetails = "SELECT COUNT(*) FROM serverdetails where domain_id = $1"

// SQLSelectServerDetails is the sql query to get server.ServerDetails struct fields in serverdetails table given a domain_id
var SQLSelectServerDetails = `
 SELECT
	id,
	address,
	ssl_grade,
	country,
	owner 
 FROM serverdetails WHERE domain_id = $1`

// SQLCheckDomainName is the sql query to get id name given a domain
var SQLCheckDomainName = `SELECT id FROM domains WHERE domain=$1;`

// SQLUpsertIntoServerDetails is the sql query to get id name given a domain
var SQLUpsertIntoServerDetails = `UPSERT INTO serverdetails 
	(id,
	address,
	ssl_grade,
	country,
	owner,
	domain_id) 
	VALUES ($1, $2, $3, $4, $5, $6);`

// SQLUpdateDomainStatus to update servers_changed & previous_ssl
var SQLUpdateDomainStatus = `UPDATE domains
SET servers_changed = $2, ssl_grade = $3, previous_ssl = $4
WHERE id = $1;`

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
)`

// SQLSelectUpdateDomain is the sql query to get update field from a given domain id
var SQLSelectUpdateDomain = `SELECT updated FROM domains WHERE id = $1`

// SQLSelectServerDetails is the sql query to get server.ServerDetails struct fields in serverdetails table given a domain_id
var SQLSelectServerDetails = `
 SELECT
	address,
	ssl_grade,
	country,
	owner 
 FROM serverdetails WHERE domain_id = $1`

// SQLCheckDomainName is the sql query to get id name given a domain
var SQLCheckDomainName = `SELECT id FROM domains WHERE domain=$1;`

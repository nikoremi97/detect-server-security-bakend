package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nikoremi97/detect-server-security-bakend/server"
	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/gofrs/uuid"
)

// DB is the database object that returns the sql.
var DB *sql.DB

// ConnectDB creates connection with database
func ConnectDB() {
	database, err := sql.Open("postgres", "postgresql://nicolas@192.168.1.58:26257/serversdb?sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	fmt.Print("CREATED DATABASE >>>")
	DB = database
	err = crdb.ExecuteTx(context.Background(), DB, nil, func(tx *sql.Tx) error {
		return setDatabase(tx)
	})
	if err == nil {
		fmt.Println("Success! database serversDB is selected")
	} else {
		log.Fatal("error: ", err)
	}
}

// setDatabase set the database serversdb in the cluster
func setDatabase(tx *sql.Tx) error {
	if _, err := tx.Exec(
		"SET database = serversDB"); err != nil {
		return err
	}
	return nil
}

// CheckDomain checks if a domain has the same name.
// if domain does not exists it calls CreateNewDomain, otherwise, it calls updateDomainFields
func CheckDomain(domain server.Domain) (server.Domain, error) {

	var domainID uuid.UUID
	row := DB.QueryRow(SQLCheckDomainName, domain.Name)
	err := row.Scan(&domainID)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")

		err = CreateNewDomain(domain)
		if err != nil {
			return domain, err
		}
	case nil:
		fmt.Println("Going to updateDomain!")
		domain, err = updateDomain(domain, domainID)
		if err != nil {
			return domain, err
		}
	default:
		fmt.Println("Erro checking domain name in database")
		if err != nil {
			return domain, err
		}
	}
	return domain, nil
}

// CreateNewDomain inserts a new domain into database
func CreateNewDomain(domain server.Domain) error {

	var newDomainID uuid.UUID
	created := time.Now().UTC()
	fmt.Println("befor DB.Exec")
	err := DB.QueryRow(SQLCreateDomain,
		domain.Name,
		domain.ServersChanged,
		domain.SslGrade,
		domain.PreviousSslGrade,
		domain.Logo,
		domain.Title,
		domain.IsDown,
		created,
	).Scan(&newDomainID)

	fmt.Println(newDomainID)
	fmt.Println("after DB.Exec")
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Error inserting new domain in DB")
	}

	fmt.Println("DOMAIN INSERTED INTO DB")

	err = createServerDetails(domain, newDomainID)
	if err != nil {
		return errors.New("Error creating ServerDetails")
	}

	fmt.Println("ALL DOMAIN DETAILS INSERTED INTO DB")
	return nil
}

// createServerDetails inserts a new domain into database
func createServerDetails(domain server.Domain, domainID uuid.UUID) error {
	var newDetailID = ""
	for _, serverDetail := range domain.Servers {

		result, err := DB.Exec(SQLCreateServerDetails,
			domainID,
			serverDetail.Address,
			serverDetail.SslGrade,
			serverDetail.Country,
			serverDetail.Owner,
		)

		fmt.Println(newDetailID)
		fmt.Println("after DB.Exec")
		if err != nil {
			fmt.Println(err.Error())
			return errors.New("Error inserting new domain in DB")

		}

		fmt.Println("DOMAIN DETAIL INSERTED INTO DB")
		fmt.Print(result)
	}

	return nil
}

// updateDomain update domain if the time diff between time.Now() and the updated field in DB for that domain is 1 minute.
// (It could be hours or whatever time metric)
func updateDomain(domain server.Domain, domainID uuid.UUID) (server.Domain, error) {

	var lastTimeUpdated time.Time
	err := DB.QueryRow(SQLGetUpdateTimeDomain, domainID).Scan(&lastTimeUpdated)

	fmt.Println(lastTimeUpdated)
	if err != nil {
		return domain, errors.New("Error geting last time updated from domain")
	}

	currentTime := time.Now().UTC()
	diff := currentTime.Sub(lastTimeUpdated)

	// convert diff to days
	minutes := int(diff.Minutes())
	if minutes > 1 {
		// with the domain is updated
		domain, err = compareServerDetails(domain, domainID)
	}

	return domain, err
}

// count the number of rows in serverdetails table where domain_id is domainID
func countServerDetails(domainID uuid.UUID) (int, error) {
	// to retrive the DetialsServers that have the domain_id = domainID, first we need to count it to create
	// an array with that length
	var count = 0
	err := DB.QueryRow(SQLCountServerDetails, domainID).Scan(&count)

	fmt.Println("in countServerDetails >>>")
	fmt.Println("count >>>")
	fmt.Print(count)

	if err != nil {
		fmt.Println("err >>>")
		fmt.Print(err)

		return 0, errors.New("Error executing query (SQLSelectServerDetails")
	}
	return count, nil
}

// count the number of rows in serverdetails table where domain_id is domainID
func countDomains() (int, error) {

	fmt.Println("here in countDomains")
	// to retrive the DetialsServers that have the domain_id = domainID, first we need to count it to create
	// an array with that length
	var count = 0
	err := DB.QueryRow(SQLCountDomains).Scan(&count)

	if err != nil {
		return 0, errors.New("Error executing query (SQLSelectServerDetails")
	}
	return count, nil
}

// getServerDetails creates an array of DetailsServer for a specific domainID.
// It also returns an array of the DetailsServer uuuid
func getServerDetails(domainID uuid.UUID) ([]server.DetailsServer, []uuid.UUID, error) {
	fmt.Println("	here in getServerDetails")
	fmt.Println(domainID)
	fmt.Println()

	var count = 0
	var err error
	count, err = countServerDetails(domainID)

	// storedServerDetails will contain the id of the serverdetial row
	storedDetailsIDs := make([]uuid.UUID, count)
	storedServerDetails := make([]server.DetailsServer, count)
	if err != nil {
		return storedServerDetails, storedDetailsIDs, errors.New("Error executing counting ServerDetails")
	}

	rows, err := DB.Query(SQLSelectServerDetails, domainID)
	if err != nil {
		return storedServerDetails, storedDetailsIDs, errors.New("Error executing query (SQLSelectServerDetails")
	}

	// to scan rows and filling storedDetailsIDs and storedServerDetails arrays
	index := 0
	defer rows.Close()
	for rows.Next() {
		var serverDetail server.DetailsServer
		err = rows.Scan(&storedDetailsIDs[index], &serverDetail.Address, &serverDetail.SslGrade, &serverDetail.Country, &serverDetail.Owner)

		if err != nil {
			return storedServerDetails, storedDetailsIDs, errors.New("Error scanning stored ServerDetails")
		}

		storedServerDetails[index] = serverDetail
		index++
	}

	fmt.Println("Stored DetailsServer retrieved")
	err = rows.Err()
	if err != nil {
		return storedServerDetails, storedDetailsIDs, err
	}

	return storedServerDetails, storedDetailsIDs, nil
}

// compareServerDetails compares domain.Servers and compare with the stored servers in DB.
// If there is a difference, update domain status fields.
func compareServerDetails(domain server.Domain, domainID uuid.UUID) (server.Domain, error) {

	storedServerDetails, storedDetailsIDs, err := getServerDetails(domainID)

	if err != nil {
		return domain, errors.New("Error executing query (SQLSelectServerDetails")
	}

	var servers = domain.Servers
	sameDetails := server.SameServerDetails(storedServerDetails, servers)
	if sameDetails {
		return domain, nil
	}

	// updateServerDetials in the DB.
	err = updateServerDetials(domainID, storedDetailsIDs, servers)
	if err != nil {
		return domain, err
	}

	// updateDomainStatus updates domain status fields in DB and return the updated domain
	domain, err = updateDomainStatus(domain, domainID)

	return domain, err
}

// updateServerDetials updates the serverDetails rows for a given domainID.
// If there is a serverDetail for a given uuid, it updates its fields.
// Otherwise, it insert another serverDetail with a new uuid and the same domainID.
func updateServerDetials(domainID uuid.UUID, storedDetailsIDs []uuid.UUID, servers []server.DetailsServer) error {

	// check if it is okay
	fmt.Println("here in updateServerDetials >>>")
	for index, serverDetail := range servers {
		var detailID uuid.UUID
		var i = index
		// if now there are more servers, it creates a new UUID
		if len(storedDetailsIDs)-i > 0 {
			detailID = storedDetailsIDs[index]
		} else {
			detailID = uuid.Must(uuid.NewV4())
		}

		fmt.Println(detailID)

		result, err := DB.Exec(SQLUpsertIntoServerDetails,
			detailID,
			serverDetail.Address,
			serverDetail.SslGrade,
			serverDetail.Country,
			serverDetail.Owner,
			domainID)

		fmt.Println(result)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateDomainStatus updates PreviousSslGrade, SslGrade, ServersChanged in database
func updateDomainStatus(domain server.Domain, domainID uuid.UUID) (server.Domain, error) {

	previousSSL, err := getPreviousSSL(domainID)
	if err != nil {
		return domain, err
	}

	domain.PreviousSslGrade = previousSSL
	domain.SslGrade = server.GetSslGrade(domain.Servers)
	domain.ServersChanged = true

	// once the domain fields are updated, it updates the domain in DB.
	_, err = DB.Exec(SQLUpdateDomainStatus, domainID, domain.ServersChanged, domain.SslGrade, domain.PreviousSslGrade)

	if err != nil {
		return domain, err

	}

	return domain, nil
}

// getPreviousSSL get the stored SSLGrade for a given domainID in the DB.
func getPreviousSSL(domainID uuid.UUID) (string, error) {

	var previousSSL = ""
	err := DB.QueryRow(SQLPreviousSSLGradeDomain, domainID).Scan(&previousSSL)

	return previousSSL, err
}

// GetDomains retrives all domains stored from database.
// Return a server.StoredDomains object which contians the domains and its serverDetails
func GetDomains() (server.StoredDomains, error) {

	fmt.Println("here in GetDomains")

	// get the correct size for array
	var totalDomains = 0
	var err error
	totalDomains, err = countDomains()

	storedDomains := server.StoredDomains{}
	storedDomains.Items = make([]server.Domain, totalDomains)

	storedDomainsIDs := make([]uuid.UUID, totalDomains)

	if err != nil {
		return storedDomains, err
	}

	// Retrives all domains table rows
	rows, err := DB.Query(SQLSelectDomains)

	if err != nil {
		return storedDomains, errors.New("Error retriving stored domains")
	}

	// Iterates over those domains rows and also create is server details
	index := 0
	defer rows.Close()
	for rows.Next() {
		fmt.Println("here reading rows")
		var domain server.Domain

		err = rows.Scan(&storedDomainsIDs[index],
			&domain.Name,
			&domain.ServersChanged,
			&domain.SslGrade,
			&domain.PreviousSslGrade,
			&domain.Logo,
			&domain.Title,
			&domain.IsDown)

		fmt.Println("domain >>> ")
		fmt.Println(domain)

		if err != nil {
			return storedDomains, errors.New("Error scanning stored domains")
		}

		// now calls getServerDetails to get domain's servers
		servers, _, err := getServerDetails(storedDomainsIDs[index])

		fmt.Println(err)

		if err != nil {
			return storedDomains, errors.New("Error retriving server details for domain")

		}

		// assign domain to items[index]
		domain.Servers = servers
		storedDomains.Items[index] = domain
		index++
	}

	return storedDomains, err
}

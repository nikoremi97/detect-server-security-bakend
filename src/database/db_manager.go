package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"../server"
	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/gofrs/uuid"
)

// DB is the database object that returns the sql.
var DB *sql.DB

// ConnectDB creates connection with database
func ConnectDB() {
	database, err := sql.Open("postgres", "postgresql://nicolas@192.168.1.58:26257/bank?sslmode=disable")
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

	var domainID = ""
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

	var newDomainID = ""
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
func createServerDetails(domain server.Domain, domainID string) error {
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

func updateDomain(domain server.Domain, domainID string) (server.Domain, error) {

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
		domain, err = compareServerDetails(domain, domainID)
	}

	return domain, err
}

// count the number of rows in serverdetails table where domain_id is domainID
func countRows(domainID string) (int, error) {
	// to retrive the DetialsServers that have the domain_id = domainID, first we need to count it to create
	// an array with that length
	var count = 0
	err := DB.QueryRow(SQLCountServerDetails, domainID).Scan(&count)

	fmt.Println("count >>>")
	fmt.Print(count)

	if err != nil {
		return 0, errors.New("Error executing query (SQLSelectServerDetails")
	}
	return count, nil
}

func compareServerDetails(domain server.Domain, domainID string) (server.Domain, error) {

	count, err := countRows(domainID)

	// storedServerDetails will cotain the id of the serverdetial row
	storedIDs := make([]uuid.UUID, count)
	storedServerDetails := make([]server.DetailsServer, count)

	rows, err := DB.Query(SQLSelectServerDetails, domainID)
	if err != nil {
		return domain, errors.New("Error executing query (SQLSelectServerDetails")
	}

	// to scan rows and filling storedIDs and storedServerDetails arrays
	index := 0
	defer rows.Close()
	for rows.Next() {
		var serverDetail server.DetailsServer
		err = rows.Scan(&storedIDs[index], &serverDetail.Address, &serverDetail.SslGrade, &serverDetail.Country, &serverDetail.Owner)

		if err != nil {
			return domain, errors.New("Error scanning saved ServerDetails from " + domainID)
		}

		storedServerDetails = append(storedServerDetails, serverDetail)
		index++
	}

	fmt.Println("Stored DetailsServer retrieved")
	err = rows.Err()
	if err != nil {
		return domain, err
	}

	var servers = domain.Servers
	sameDetails := server.SameServerDetails(storedServerDetails, servers)
	if sameDetails {
		return domain, nil
	}

	err = updateServerDetials(domainID, storedIDs, servers)
	if err != nil {
		return domain, err
	}

	domain, err = updateDomainStatus(domain, domainID)

	return domain, err
}

func updateServerDetials(domainID string, storedIDs []uuid.UUID, servers []server.DetailsServer) error {

	fmt.Println("here in updateServerDetials >>>")
	for index, serverDetail := range servers {
		var detailID uuid.UUID
		if ((len(storedIDs) - 1) - 1) > 0 {
			detailID = storedIDs[index]
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

func updateDomainStatus(domain server.Domain, domainID string) (server.Domain, error) {

	previousSSL, err := getPreviousSSL(domainID)
	if err != nil {
		return domain, err
	}

	domain.PreviousSslGrade = previousSSL
	domain.SslGrade = server.GetSslGrade(domain.Servers)
	domain.ServersChanged = true
	_, err = DB.Exec(SQLUpdateDomainStatus, domainID, domain.ServersChanged, domain.SslGrade, domain.PreviousSslGrade)

	if err != nil {
		return domain, err

	}

	return domain, nil
}

func getPreviousSSL(domainID string) (string, error) {

	var previousSSL = ""
	err := DB.QueryRow(SQLPreviousSSLGradeDomain, domainID).Scan(&previousSSL)

	return previousSSL, err
}

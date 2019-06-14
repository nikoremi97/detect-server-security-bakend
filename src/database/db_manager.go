package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cockroachdb/cockroach-go/crdb"

	"../server"
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
	err := DB.QueryRow(SQLSelectUpdateDomain, domainID).Scan(&lastTimeUpdated)

	fmt.Println(lastTimeUpdated)
	if err != nil {
		fmt.Println("shit men")
		return domain, errors.New("Error geting last time updated from domain")
	}

	currentTime := time.Now().UTC()
	diff := currentTime.Sub(lastTimeUpdated)

	// convert diff to days
	hours := int(diff.Hours())
	if hours > 1 {
		compareServerDetails(domain.Servers, domainID)
	}

	return domain, err
}

func compareServerDetails(servers []server.DetailsServer, domainID string) (bool, error) {

	rows, err := DB.Query(SQLSelectServerDetails, domainID)

	if err != nil {
		return false, errors.New("Error executing query (SQLSelectServerDetails")
	}

	savedServerDetails := []server.DetailsServer{}
	defer rows.Close()
	for rows.Next() {
		var serverDetail server.DetailsServer
		err = rows.Scan(&serverDetail.Address, &serverDetail.SslGrade, &serverDetail.Country, &serverDetail.Owner)
		if err != nil {
			return false, errors.New("Error scanning saved ServerDetails from " + domainID)
		}
		savedServerDetails = append(savedServerDetails, serverDetail)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return false, err
	}

	sameDetails := server.SameServerDetails(savedServerDetails, servers)
	if sameDetails {
		return false, nil
	}

	return false, err
}

# detect-server-security-bakend

detect-server-security-bakend is the backend of SecurityServer App https://github.com/nikoremi97/detect-server-security-frontend-app.

It was made using the following technologies:
- Go as language
- Chi as the API router
- CockroachDB as its DB

To run this project, do the following steps:

1. Clone this project

2. Install Go in your machine -> https://golang.org/doc/install

3. Add these packages to Go by running go get -d -v ./...
3.1 `github.com/asaskevich/govalidator`  
3.2 `github.com/likexian/whois-go`  
3.3 `github.com/go-chi/chi`  
3.4 `github.com/cockroachdb/cockroach-go/crdb`  
3.5	`github.com/gofrs/uuid`  

4. Install CockroachDB in your machine -> https://www.cockroachlabs.com/docs/stable/install-cockroachdb-linux.html

5. Init a local cluster:  
5.1  Follow this guide changing the ip address as `192.168.1.58` instead of localhost -> https://www.cockroachlabs.com/docs/stable/start-a-local-cluster.html.  
5.2 Start the sql cmd with the command `cockroach sql --insecure --host=192.168.1.58:26257`  
5.3 To set up the DB, run the sql lines in `src/database/db_schema.sql` separated by `;`  

6. On path `src/main` run `go run main.go`

7. To generate Go binary run this command on linux:
    `CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .`

Enjoy!
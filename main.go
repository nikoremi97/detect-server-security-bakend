package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nikoremi97/detect-server-security-bakend/database"
	"github.com/nikoremi97/detect-server-security-bakend/handler"
)

// initRouter creates a new Chi Router and defines its routes and methods
func initRouter() {

	r := chi.NewRouter()
	r.Route("/detectServerSecurity/api/v1/", func(r chi.Router) {
		r.Get("/getServers", handler.ExecuteGetRequest())
		r.Post("/newServer", handler.ExecutePostRequest())

	})

	// Backend listen and serve in 192.168.1.57 port 3000
	address := "192.168.1.57:3000"
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting server for testing HTTP...\n")
}

// this is the main function of the backend. It creates the connection to Database, the Chi API router and defines its routes.
func main() {go 
	database.ConnectDB()
	initRouter()
}

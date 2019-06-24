package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"../database"
	"../handler"
)

func initRouter() {

	r := chi.NewRouter()
	r.Route("/detectServerSecurity/api/v1/", func(r chi.Router) {
		r.Get("/getServers", handler.ExecuteGetRequest())
		r.Post("/newServer", handler.ExecutePostRequest())

	})

	fmt.Printf("Starting server for testing HTTP...\n")

	// Backend listen and serve in 192.168.1.57 port 3000
	address := "192.168.1.57:3000"
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(err)
	}
}

// this is the main function of the backend. It creates the connection to Database, the Chi API router and defines its routes.
func main() {
	database.ConnectDB()
	initRouter()
}

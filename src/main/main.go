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
	servers := []string{"hello server1", "hello server2", "hello server3"}

	r := chi.NewRouter()
	r.Route("/detectServerSecurity/api/v1/", func(r chi.Router) {
		r.Get("/serverStatus", handler.ExecuteGetRequest(servers))
		r.Post("/newServer", handler.ExecutePostRequest())

	})

	fmt.Printf("Starting server for testing HTTP...\n")

	// Backend listen and serve in port 3000
	port := ":3000"
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

// this is the main function of the backend. It creates the connection to Database, the Chi API router and defines its routes.
func main() {
	database.ConnectDB()
	initRouter()
}

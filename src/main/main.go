package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"../handler"
)

// this is the main function of the backend. It creates the Chi API router and defines its routes
func main() {
	servers := []string{"hello server1", "hello server2", "hello server3"}

	r := chi.NewRouter()
	r.Route("/detectServerSecurity/api/v1/", func(r chi.Router) {
		r.Get("/serverStatus", handler.ServerStatusGet(servers))
		r.Post("/newServer", handler.ServerStatusPost())

	})

	fmt.Printf("Starting server for testing HTTP...\n")

	// Backend listen and serve in port 3000
	port := ":3000"
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

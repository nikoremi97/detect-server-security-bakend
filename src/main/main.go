package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"../handler"
)

func main() {
	port := ":3000"
	servers := []string{"hello server1", "hello server2", "hello server3"}

	r := chi.NewRouter()
	r.Route("/detectServerSecurity/api/v1/", func(r chi.Router) {
		r.Get("/serverStatus", handler.ServerStatusGet(servers))
		r.Post("/newServer", handler.ServerStatusPost(""))

	})

	fmt.Printf("Starting server for testing HTTP...\n")

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

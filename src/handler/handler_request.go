package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "../database"
	server "../server"
	utils "../utils"
)

// ExecuteGetRequest get method
func ExecuteGetRequest(response []string) http.HandlerFunc {
	fmt.Println(response)

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}
}

// ExecutePostRequest get user typed domain and create a Server struct
func ExecutePostRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get form of request
		if err := r.ParseForm(); err != nil {
			fmt.Println(w, "ParseForm() err: %v", err)

			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// get domain field of request
		var domainToAnalyze = ""
		domainToAnalyze = r.FormValue("domain")

		// domainToAnalyze = strings.ToLower(domainToAnalyze)
		domainQuery, err := utils.ValidateQuery(domainToAnalyze)

		// handle if there is an error in domain typed by user
		if err != nil {

			// if the domain is not valid, stop and send StatusBadRequest
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Analyze domain and creates a new Domain struct
		newDomain, err := server.AnalyzeDomain(domainQuery)

		// handle if there is an error in domain typed by user
		if err != nil {
			// if the domain is not valid, stop and send StatusBadRequest

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newDomain, err = database.CheckDomain(newDomain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		fmt.Println("")
		fmt.Println("serversDescription >>>>")
		fmt.Print(newDomain)
		json.NewEncoder(w).Encode(newDomain)
		w.WriteHeader(http.StatusOK)

	}
}

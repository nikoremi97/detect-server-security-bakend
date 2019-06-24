package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	database "../database"
	server "../server"
	utils "../utils"
)

// ExecuteGetRequest get method
func ExecuteGetRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		storedDomains, err := database.GetDomains()

		if err != nil {
			json.NewEncoder(w).Encode(err)
			w.WriteHeader(http.StatusInternalServerError)

		}

		fmt.Print(storedDomains)
		json.NewEncoder(w).Encode(storedDomains)
		w.WriteHeader(http.StatusOK)
	}
}

// ExecutePostRequest get user typed domain and create a Server struct
func ExecutePostRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the data in the request's body and creates a []byte called body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(w, "Error reading request body", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
		}

		var t = &server.BodyRequest{}
		err = json.Unmarshal(body, t)
		if err != nil {
			fmt.Println(w, "Error Unmarshall bodyRequest", err)
			http.Error(w, "Error Unmarshall bodyRequest", http.StatusInternalServerError)
		}

		// get domain field of request
		var domainToAnalyze = t.DomainRequestParam
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
		fmt.Println(newDomain)
		json.NewEncoder(w).Encode(newDomain)
		w.WriteHeader(http.StatusOK)

	}
}

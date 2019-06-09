package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	server "../server"
	utils "../utils"
)

// ServerStatusGet get method
func ServerStatusGet(response []string) http.HandlerFunc {
	fmt.Println(response)

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}
}

var analyzerEndpoint = "https://api.ssllabs.com/api/v3/analyze?host="

// ServerStatusPost get user typed domain and create a Server struct
func ServerStatusPost() http.HandlerFunc {
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

		// make request to test server until its response contains endpoints field
		var endpointsFound = false
		for !endpointsFound {

			// make get request to Test Server analyzerEndpoint
			response, err := http.Get(analyzerEndpoint + domainQuery)
			if err != nil {

				// if there is an error in request, stop and send StatusInternalServerError
				http.Error(w, "Error trying to reach test server", http.StatusInternalServerError)
				return
			}

			// Read body form response from Test Server analyzer
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				// if there is an error in reading body, stop and send StatusInternalServerError
				http.Error(w, "Error reading body of response from test server", http.StatusInternalServerError)
				return

			}

			// convert data into json with Domain struct
			var domainInfo server.Domain
			err = json.Unmarshal(data, &domainInfo)
			if err != nil {
				http.Error(w, "Unmarshal had failed", http.StatusInternalServerError)
				return
			}

			// response from test server could be OK, but domain status could be error.
			// the most common case is when the domain is a valid DNSName but the is no a real host with
			// that dns name
			if domainInfo.Status == "ERROR" {
				endpointsFound = true

				http.Error(w, "Error in test server. Try with another domain", http.StatusNotFound)
				return

			}

			// when domain contains Endpoints, create Server struct and send OK status
			if len(domainInfo.Endpoints) != 0 {

				endpointsFound = true

				var server = server.CreateServer(domainInfo)

				fmt.Println("")
				fmt.Println("serversDescription >>>>")
				json.NewEncoder(w).Encode(server)
				w.WriteHeader(http.StatusOK)

			}

		}
	}
}

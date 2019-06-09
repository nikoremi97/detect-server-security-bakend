package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	server "../server"
)

/**
* ServerStatusGet get method
 */
func ServerStatusGet(response []string) http.HandlerFunc {
	fmt.Println(response)

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}
}

var analyzerEndpoint = "https://api.ssllabs.com/api/v3/analyze?host="

// ServerStatusPost get user typed domain and create a Server struct
func ServerStatusPost(domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HERE IN POST")
		var domainToAnalyze = domain
		fmt.Println("DOMAIN>>> " + domainToAnalyze)
		if err := r.ParseForm(); err != nil {
			fmt.Println(w, "ParseForm() err: %v", err)
			return
		}

		if domainToAnalyze == "" {
			domainToAnalyze = r.FormValue("domain")
			fmt.Println("DOMAIN from request>>> " + domainToAnalyze)
		}

		response, err := http.Get(analyzerEndpoint + domainToAnalyze)
		if err != nil {
			fmt.Println("The HTTP request failed with error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading body of response")
			w.WriteHeader(http.StatusInternalServerError)
		}

		var domainInfo server.Domain
		err = json.Unmarshal(data, &domainInfo)
		if err != nil {
			fmt.Println("Unmarshal IS FAILED", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		var buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.Encode(domainInfo)

		if domainInfo.Status == "ERROR" {

			w.WriteHeader(http.StatusBadRequest)
		} else {

			if len(domainInfo.Endpoints) != 0 {
				var server = server.CreateServer(domainInfo)

				fmt.Println("")
				fmt.Println("serversDescription >>>>")
				json.NewEncoder(w).Encode(server)
				w.WriteHeader(http.StatusOK)

			} else {
				fmt.Println("retrying reaching endpoints >>>")

				go ServerStatusPost(domainToAnalyze)
			}
		}
	}
}

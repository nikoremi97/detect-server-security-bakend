package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
			fmt.Fprintf(w, "domain = %s\n", domainToAnalyze)
			fmt.Println("DOMAIN from request>>> " + domainToAnalyze)
		}

		response, err := http.Get(analyzerEndpoint + domainToAnalyze)
		if err != nil {

			log.Fatalln("The HTTP request failed with error", err)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {

			log.Fatalln("Error reading body of response")
		}

		// fmt.Println(string(data))

		var domainInfo server.Domain
		err = json.Unmarshal(data, &domainInfo)
		if err != nil {
			log.Fatalln("Unmarshal IS FAILED", err)
		}

		var buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.Encode(domainInfo)

		fmt.Println("domain info >>> ")
		fmt.Println(domainInfo)

		fmt.Println("domainInfo.Endpoints >>> ")
		fmt.Println(domainInfo.Endpoints)

		if len(domainInfo.Endpoints) != 0 {

			var server = server.CreateServer(domainInfo)
			// var serversDescription = server.CreateServersDescripton(domainInfo.Endpoints)

			fmt.Println("")
			fmt.Println("serversDescription >>>>")

			var buf = new(bytes.Buffer)
			enc := json.NewEncoder(buf)
			enc.Encode(server)
			// json.NewEncoder(w).Encode(serversDescription)
			fmt.Println(buf)
		} else {
			fmt.Println("retrying reaching endpoints >>>")

			ServerStatusPost(domainToAnalyze)
		}

		w.WriteHeader(http.StatusOK)
	}
}

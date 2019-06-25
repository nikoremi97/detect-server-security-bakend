package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var enpointAnalyzer = "https://api.ssllabs.com/api/v3/analyze?host="

// AnalyzeDomain make requests to SSLLABS analyze endpoint and returns the created Domain based on it
func AnalyzeDomain(domainQuery string) (Domain, error) {

	newDomain := Domain{}
	// make request to test server until its response contains endpoints field
	var endpointsFound = false
	for !endpointsFound {

		// make get request to Test Server enpointAnalyzer
		response, err := http.Get(enpointAnalyzer + domainQuery)
		if err != nil {

			// if there is an error in request, stop and send StatusInternalServerError
			return newDomain, errors.New("Error trying to reach analyzer server")
		}

		if response.StatusCode != http.StatusOK {
			return newDomain, errors.New("Error in analyzer server")

		}

		// Read body form response from Test Server analyzer
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			// if there is an error in reading body, stop and send StatusInternalServerError
			return newDomain, errors.New("Error reading body of response from test server")

		}

		// convert data into json with Domain struct
		var domainInfo DomainDescription
		err = json.Unmarshal(data, &domainInfo)
		if err != nil {
			return newDomain, errors.New("Unmarshal had failed")
		}

		// response from test server could be OK, but domain status could be error.
		// the most common case is when the domain is a valid DNSName but the is no a real host with
		// that dns name
		if domainInfo.Status == "ERROR" {
			endpointsFound = true

			return newDomain, errors.New("Error in test server. Try with another domain")
		}

		// when domain contains Endpoints, create Server struct and send OK status
		if len(domainInfo.Endpoints) != 0 {

			endpointsFound = true

			newDomain, err = CreateServer(domainInfo)

			if err != nil {
				return newDomain, errors.New("Creating new Server")

			}
		}
	}

	return newDomain, nil
}

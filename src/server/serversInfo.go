package server

// Imports
import (
	"fmt"
	"strings"

	"github.com/likexian/whois-go"
)

const ipv6Address = "CAN´T ANALYZE IPV6 ADDRESS"
const unknown = "UNKNOWN"

// SSLGrades is an array with the posible ssl grade values
var SSLGrades = []string{"A+", "A", "B", "C", "D", "E", "F", "T", "M"}

// CreateServer classifies endpoints into DescriptionServer structures.
func CreateServer(domain Domain) Server {
	var server = Server{
		Servers:          nil,
		Title:            "title",
		Logo:             "logo",
		ServersChanged:   false,
		SslGrade:         "T",
		PreviousSslGrade: "T",
		IsDown:           false,
	}

	server = ConfigureServerDescription(domain, server)
	return server
}

// ConfigureServerDescription classifies endpoints into DescriptionServer structures.
func ConfigureServerDescription(domain Domain, server Server) Server {
	var newServer = server
	var endpoints = domain.Endpoints
	var serversDescription = make([]DescriptionServer, len(endpoints))

	newServer.Servers = serversDescription

	for i := 0; i < len(endpoints); i++ {

		serversDescription[i].SslGrade = endpoints[i].Grade
		serversDescription[i].Address = endpoints[i].IPAddress

		if strings.Contains(endpoints[i].IPAddress, ":") {
			fmt.Println("IS IPV6 >>> ")

			serversDescription[i].Country = ipv6Address
			serversDescription[i].Owner = ipv6Address
			continue
		}

		fmt.Println("WHO IS >>>>")
		result, err := whois.Whois(endpoints[i].IPAddress)
		if err != nil {
			fmt.Println("WHOIS COMMAND FAILED")
			newServer.IsDown = true

			serversDescription[i].Country = unknown
			serversDescription[i].Owner = unknown
			continue
		}

		fmt.Println("WHO IS SUCCESS >>>>")
		serversDescription[i] = getOwnerAndCountry(serversDescription[i], result)

	}
	newServer.SslGrade = getSslGrade(newServer.Servers)
	newServer.PreviousSslGrade = getPreviousSSL(newServer.Servers)

	return newServer
}

// getOwner
func getOwnerAndCountry(descriptionServer DescriptionServer, result string) DescriptionServer {
	for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {

		hasOrgName := strings.HasPrefix(line, "OrgName")
		hasCountry := strings.HasPrefix(line, "Country")
		if hasOrgName || hasCountry {
			fmt.Println("ORG >>>" + line)

			infoArray := strings.Split(line, ":")

			info := strings.TrimSpace(infoArray[1])

			if hasOrgName {
				descriptionServer.Owner = info
			} else if hasCountry {
				descriptionServer.Country = info
			}
		}
	}
	return descriptionServer
}

// getSslGrade gets the lower grade of the SSLGrade in Server endpoints
func getSslGrade(descriptionServer []DescriptionServer) string {
	var sslGrade = ""
	var gradeIndex = -1

	for _, endpoint := range descriptionServer {

		var currentGrade = indexOf(endpoint.SslGrade, SSLGrades)
		if currentGrade > gradeIndex {
			gradeIndex = currentGrade

		}
	}

	sslGrade = SSLGrades[gradeIndex]
	return sslGrade
}

// getSslGrade gets the lower grade of the SSLGrade in Server endpoints
func getPreviousSSL(descriptionServer []DescriptionServer) string {
	return "TODO"
}

// indexOf find index of element in data array
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

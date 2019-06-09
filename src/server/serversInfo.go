package server

// Imports
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	utils "../utils"
	"github.com/likexian/whois-go"
)

const ipv6AddressWarning = "CAN´T ANALYZE IPV6 ADDRESS"
const unknown = "UNKNOWN"

// SSLGrades is an array with the posible ssl grade values
var SSLGrades = []string{"A+", "A", "B", "C", "D", "E", "F", "T", "M"}

// CreateServer classifies endpoints into DescriptionServer structures.
func CreateServer(domain Domain) Server {

	infoTitleLogo := getTitleOrLogo(domain.Host)

	var server = Server{
		Servers:          nil,
		Title:            infoTitleLogo[0],
		Logo:             infoTitleLogo[1],
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

		if endpoints[i].StatusMessage == "Ready" {
			serversDescription[i].SslGrade = endpoints[i].Grade

		} else {
			serversDescription[i].SslGrade = unknown

		}

		serversDescription[i].Address = endpoints[i].IPAddress

		if strings.Contains(endpoints[i].IPAddress, ":") {
			fmt.Println("IS IPV6 >>> ")

			serversDescription[i].Country = ipv6AddressWarning
			serversDescription[i].Owner = ipv6AddressWarning
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
	fmt.Println("here in getOwnerAndCountry >>>>>>>>> ")

	fmt.Println(descriptionServer)

	var arinWhoIs = "whois.arin.net"
	var apnicWhoIs = "whois.apnic.net"
	var isArinHost = false
	var isApnicHost = false
	var hasOrgName = false
	var hasCountry = false
	var info = ""

	for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {

		whoIsRefer := strings.HasPrefix(line, "whois")
		if whoIsRefer {
			infoWhoIsRefer := strings.Split(line, ":")

			whoIs := strings.TrimSpace(infoWhoIsRefer[1])

			if whoIs == arinWhoIs {
				isArinHost = true
			} else if whoIs == apnicWhoIs {
				isApnicHost = true
			}
			continue
		}

		if isArinHost {
			hasOrgName = strings.HasPrefix(line, "OrgName")
			hasCountry = strings.HasPrefix(line, "Country")

		} else if isApnicHost {
			hasOrgName = strings.HasPrefix(line, "descr")
			hasCountry = strings.HasPrefix(line, "country")

		}

		if hasOrgName || hasCountry {
			infoArray := strings.Split(line, ":")
			info = strings.TrimSpace(infoArray[1])

			if hasOrgName && descriptionServer.Owner == "" {
				descriptionServer.Owner = info
			} else if hasCountry {
				descriptionServer.Country = info
				break
			}
		}
	}
	return descriptionServer
}

// getSslGrade gets the lower grade of the SSLGrade in Server endpoints
func getSslGrade(descriptionServer []DescriptionServer) string {

	fmt.Println("HERE IN getSslGrade")
	fmt.Println(descriptionServer)
	var sslGrade = ""
	var gradeIndex = -1

	for _, endpoint := range descriptionServer {
		fmt.Println(endpoint.SslGrade)
		if endpoint.SslGrade != unknown {

			var currentGrade = utils.IndexOf(endpoint.SslGrade, SSLGrades)
			if currentGrade > gradeIndex {
				gradeIndex = currentGrade

			}
			sslGrade = SSLGrades[gradeIndex]
		}

		if sslGrade == "" {

			sslGrade = endpoint.SslGrade
		}

	}

	return sslGrade
}

// getSslGrade gets the lower grade of the SSLGrade in Server endpoints
func getPreviousSSL(descriptionServer []DescriptionServer) string {
	return "TODO"
}

// getTitle from the head of the host webpage
func getTitleOrLogo(hostName string) []string {
	var title = ""
	var logo = ""
	var info = []string{"ERROR", "ERROR"}
	response, err := http.Get("https://" + hostName)
	if err != nil {
		fmt.Println("The HTTP request failed with error", err)

		return info
	}

	if response.StatusCode == http.StatusOK {

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading body of response")
			return info
		}

		content := string(data)

		for _, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {

			isTitle := strings.Contains(line, "<title>")
			if isTitle && title == "" {
				title = utils.TrimTitle(line)
				break
			}
		}

		for _, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {
			isLogo := strings.Contains(line, `type="image/x-icon"`)
			if isLogo && logo == "" {
				logo = utils.TrimLogo(line)
				break

			}
		}

	}

	info[0] = title
	info[1] = logo
	return info
}

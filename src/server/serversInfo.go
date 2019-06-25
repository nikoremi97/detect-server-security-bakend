package server

// Imports
import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	utils "../utils"
	"github.com/likexian/whois-go"
)

// CreateServer classifies endpoints into Domain structures.
func CreateServer(domainDescription DomainDescription) (Domain, error) {
	fmt.Println("here in CreateServer")

	infoTitleLogo := getTitleOrLogo(domainDescription.Host)

	domain := Domain{
		Name:             domainDescription.Host,
		Servers:          nil,
		Title:            infoTitleLogo[0],
		Logo:             infoTitleLogo[1],
		ServersChanged:   false,
		SslGrade:         "T",
		PreviousSslGrade: "T",
		IsDown:           false,
	}

	domainFilled, err := ConfigureServerDescription(domainDescription, domain)
	return domainFilled, err
}

// ConfigureServerDescription classifies endpoints into DetailsServer structures.
func ConfigureServerDescription(domainDescription DomainDescription, domain Domain) (Domain, error) {
	fmt.Println("here in ConfigureServerDescription")
	newDomain := Domain{}
	var endpoints = domainDescription.Endpoints
	var detailsServer = make([]DetailsServer, len(endpoints))

	domain.Servers = detailsServer

	for i := 0; i < len(endpoints); i++ {

		if endpoints[i].StatusMessage == "Ready" {
			detailsServer[i].SslGrade = endpoints[i].Grade

		} else {
			detailsServer[i].SslGrade = UNKNOWN

		}

		detailsServer[i].Address = endpoints[i].IPAddress

		if strings.Contains(endpoints[i].IPAddress, ":") {
			fmt.Println("IS IPV6 >>> ")

			detailsServer[i].Country = Ipv6AddressWarning
			detailsServer[i].Owner = Ipv6AddressWarning
			continue
		}

		fmt.Println("WHO IS >>>>")
		result, err := whois.Whois(endpoints[i].IPAddress)
		if err != nil {
			fmt.Println("WHOIS COMMAND FAILED")
			domain.IsDown = true

			detailsServer[i].Country = UNKNOWN
			detailsServer[i].Owner = UNKNOWN

			return newDomain, errors.New("WhoIs command failed")
		}

		fmt.Println("WHO IS SUCCESS >>>>")

		detailsServer[i] = getOwnerAndCountry(detailsServer[i], result)

	}
	domain.SslGrade = GetSslGrade(domain.Servers)
	// because the first time the previous sslgrade is the same as the current one
	domain.PreviousSslGrade = domain.SslGrade

	newDomain = domain
	return newDomain, nil
}

// getOwner
func getOwnerAndCountry(descriptionServer DetailsServer, result string) DetailsServer {
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

// GetSslGrade gets the lower grade of the SSLGrade in Server endpoints
func GetSslGrade(descriptionServer []DetailsServer) string {

	fmt.Println("HERE IN GetSslGrade")
	fmt.Println(descriptionServer)
	var sslGrade = ""
	var gradeIndex = -1

	for _, endpoint := range descriptionServer {
		fmt.Println(endpoint.SslGrade)

		if endpoint.SslGrade != UNKNOWN {
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

// GetPreviousSSL gets the previous SSL grade in Domain
func GetPreviousSSL() string {

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

// SameServerDetails returns wheather the domain has still the same DetailsServer or not.
func SameServerDetails(oldDetailsServer []DetailsServer, currentDetailsServer []DetailsServer) bool {

	// First lets ensure both have the same length
	if len(oldDetailsServer) != len(currentDetailsServer) {
		return false
	}

	for index, oldDetail := range oldDetailsServer {
		currentDetail := currentDetailsServer[index]

		if (oldDetail.Address != currentDetail.Address) ||
			(oldDetail.Country != currentDetail.Country) ||
			(oldDetail.Owner != currentDetail.Owner) ||
			(oldDetail.SslGrade != currentDetail.SslGrade) {
			return false

		}
	}

	return true
}

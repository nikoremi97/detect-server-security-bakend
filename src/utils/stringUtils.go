package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
)

// IndexOf find index of element in data array
func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

// ValidateQuery validates if query has a valid domain
func ValidateQuery(queryParam string) (fixedQuery string, err error) {
	fmt.Println("here in ValidateQuery")
	fmt.Println(queryParam)
	var query = queryParam
	query = strings.TrimSpace(query)
	query = strings.ToLower(query)

	spacesInQuery := strings.Split(query, " ")
	if len(spacesInQuery) > 1 {
		return "", errors.New("domain name must not contain spaces")
	}

	dotsInQuery := strings.Split(query, ".")
	if len(dotsInQuery) < 2 {
		return "", errors.New("domain name must contain at least one dot")
	}

	isDNSName := govalidator.IsDNSName(query)
	if isDNSName {
		return query, nil

	}

	return "", errors.New("You typed and invalid domain name")
}

// TrimTitle takes the like that contians title tag and split its content to get title
func TrimTitle(line string) string {
	var title = ""

	line = strings.TrimSpace(line)

	// <title>Example</title>
	splitedTitle := strings.Split(line, "<title>")

	// Example</title>
	title = splitedTitle[1]

	splitedTitle = strings.Split(title, "</title>")

	// Example
	title = splitedTitle[0]
	return title
}

// TrimLogo search the tag type="image/x-icon", where its contain page's logo.
func TrimLogo(line string) string {
	var logo = ""
	var lineInfo = ""
	line = strings.TrimSpace(line)

	// split the type="image/x-icon" section
	splitedlogo := strings.Split(line, `type="image/x-icon"`)

	// take the first element of array
	lineInfo = splitedlogo[0]

	// split the lineInfo info in <link
	splitedlogo = strings.Split(lineInfo, "<link")

	// take the last element of splited lineInfo because the logo link is the closest element from type="image/x-icon"
	maxPosition := len(splitedlogo)
	if maxPosition > 0 {
		lineInfo = splitedlogo[maxPosition-1]

	} else {
		lineInfo = splitedlogo[0]
	}

	// split the href tag
	splitedlogo = strings.Split(lineInfo, "href=")

	for _, linkLogo := range splitedlogo {

		linkLogo = strings.TrimSpace(linkLogo)
		// link starts with `"` character
		isLinkLogo := strings.HasPrefix(linkLogo, `"`)
		if isLinkLogo {

			// fixing unwanted added characters
			linkLogo = strings.Trim(linkLogo, `"`)

			// it is posible that in the same line that is the logo, another info could be present
			// so, it has to be splitted by " "
			splitedlogo = strings.Split(linkLogo, " ")
			linkLogo = splitedlogo[0]
			lineInfo = strings.Trim(linkLogo, `"`)

			// to avoid wrong logo links, the logo var is only assigned when all the process is completed
			logo = lineInfo
		}
	}
	return logo
}

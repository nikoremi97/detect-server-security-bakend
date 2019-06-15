package server

// CONSTANTS

// Ipv6AddressWarning is the warning when the IP is ipv6
const Ipv6AddressWarning = "CAN´T ANALYZE IPV6 ADDRESS"

// UNKNOWN is the value when can´t analyze the server details
const UNKNOWN = "UNKNOWN"

// SSLGrades is an array with the posible ssl grade values
var SSLGrades = []string{"A+", "A", "B", "C", "D", "E", "F", "T", "M"}

// Structure types

// DomainDescription is the struct of the SSL Test report
type DomainDescription struct {
	Host            string                `json:"host,omitempty"`
	Port            int                   `json:"port,omitempty"`
	Protocol        string                `json:"protocol,omitempty"`
	IsPublic        bool                  `json:"isPublic,omitempty"`
	Status          string                `json:"status,omitempty"`
	StartTime       int64                 `json:"startTime,omitempty"`
	TestTime        int64                 `json:"testTime,omitempty"`
	EngineVersion   string                `json:"engineVersion,omitempty"`
	CriteriaVersion string                `json:"criteriaVersion,omitempty"`
	Endpoints       []EndpointDescription `json:"endpoints,omitempty"`
}

// EndpointDescription is the struct of an Endpoint in the SSL Test report
type EndpointDescription struct {
	IPAddress         string `json:"ipAddress,omitempty"`
	ServerName        string `json:"serverName,omitempty"`
	StatusMessage     string `json:"statusMessage,omitempty"`
	Grade             string `json:"grade,omitempty"`
	GradeTrustIgnored string `json:"gradeTrustIgnored,omitempty"`
	HasWarnings       bool   `json:"hasWarnings,omitempty"`
	IsExceptional     bool   `json:"isExceptional,omitempty"`
	Progress          int    `json:"progress,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	Delegation        int    `json:"delegation,omitempty"`
}

// Domain is the struct of the data stored in the DB
type Domain struct {
	Name             string          `json:"domain"`
	Servers          []DetailsServer `json:"servers"`
	ServersChanged   bool            `json:"servers_changed"`
	SslGrade         string          `json:"ssl_grade"`
	PreviousSslGrade string          `json:"previous_ssl_grade"`
	Logo             string          `json:"logo"`
	Title            string          `json:"title"`
	IsDown           bool            `json:"is_down"`
}

// DetailsServer is the struct of a server in Server
type DetailsServer struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

// StoredDomains is the struct t
type StoredDomains struct {
	Items []Domain `json:"items"`
}

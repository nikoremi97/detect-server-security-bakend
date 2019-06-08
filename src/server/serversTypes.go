package server

// Domain is the struct of the SSL Test report
type Domain struct {
	Host            string     `json:"host,omitempty"`
	Port            int        `json:"port,omitempty"`
	Protocol        string     `json:"protocol,omitempty"`
	IsPublic        bool       `json:"isPublic,omitempty"`
	Status          string     `json:"status,omitempty"`
	StartTime       int64      `json:"startTime,omitempty"`
	TestTime        int64      `json:"testTime,omitempty"`
	EngineVersion   string     `json:"engineVersion,omitempty"`
	CriteriaVersion string     `json:"criteriaVersion,omitempty"`
	Endpoints       []Endpoint `json:"endpoints,omitempty"`
}

// Endpoint is the struct of an Endpoint in the SSL Test report
type Endpoint struct {
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

// Server is the struct of the data stored in the DB
type Server struct {
	Servers          []DescriptionServer `json:"servers"`
	ServersChanged   bool                `json:"servers_changed"`
	SslGrade         string              `json:"ssl_grade"`
	PreviousSslGrade string              `json:"previous_ssl_grade"`
	Logo             string              `json:"logo"`
	Title            string              `json:"title"`
	IsDown           bool                `json:"is_down"`
}

// DescriptionServer is the struct of a server in Server
type DescriptionServer struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

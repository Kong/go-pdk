package entities

type ServiceKey struct {
	ID string `json:"id"`
}

type CertificateKey struct {
	ID string `json:"id"`
}

type Route struct {
	ID                      string      `json:"id"`
	CreatedAt               int         `json:"created_at"`
	UpdatedAt               int         `json:"updated_at"`
	Name                    string      `json:"name"`
	Protocols               []string    `json:"protocols"`
	Methods                 *[]string   `json:"methods"`
	Hosts                   *[]string   `json:"hosts"`
	Paths                   *[]string   `json:"paths"`
	Headers                 *[]string   `json:"headers"`
	HTTPSRedirectStatusCode int         `json:"https_redirect_status_code"`
	RegexPriority           int         `json:"regex_priority"`
	StripPath               bool        `json:"strip_path"`
	PreserveHost            bool        `json:"preserve_host"`
	SNIs                    *[]string   `json:"snis"`
	Sources                 *[]string   `json:"sources"`
	Destinations            *[]string   `json:"destinations"`
	Tags                    *[]string   `json:"tags"`
	Service                 *ServiceKey `json:"service"`
}

type Service struct {
	ID                string          `json:"id"`
	CreatedAt         int             `json:"created_at"`
	UpdatedAt         int             `json:"updated_at"`
	Name              string          `json:"name"`
	Retries           int             `json:"retries"`
	Protocol          string          `json:"protocol"`
	Host              string          `json:"host"`
	Port              int             `json:"port"`
	Path              string          `json:"path"`
	ConnectTimeout    int             `json:"connect_timeout"`
	WriteTimeout      int             `json:"write_timeout"`
	ReadTimeout       int             `json:"read_timeout"`
	Tags              *[]string       `json:"tags"`
	ClientCertificate *CertificateKey `json:"client_certificate"`
}

type Consumer struct {
	ID        string    `json:"id"`
	CreatedAt int       `json:"created_at"`
	Username  *string   `json:"username"`
	CustomID  *string   `json:"custom_id"`
	Tags      *[]string `json:"tags"`
}

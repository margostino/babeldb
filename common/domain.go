package common

const (
	SuccessEnabled = iota
	FailureEnabled
	RandomnessEnabled
)

type Striker struct {
	Url string `yaml:"url"`
}

type Client struct {
	Url         string `yaml:"url"`
	RequestFile string `yaml:"requestFile"`
	CallsNumber int    `yaml:"callsNumber"`
	MaxStepTime int    `yaml:"maxStepTime"`
	Method      string `yaml:"method"`
	ContentType string `yaml:"contentType"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Enabled     bool   `yaml:"enabled"`
}

type Server struct {
	Port            string `yaml:"port"`
	Host            string `yaml:"host"`
	Path            string `yaml:"path"`
	ResponseFile    string `yaml:"responseFile"`
	HealthcheckPath string `yaml:"healthcheckPath"`
	HealthcheckFile string `yaml:"healthcheckFile"`
	HotStatusPath   string `yaml:"hotStatusPath"`
}

type Configuration struct {
	Clients []*Client `yaml:"clients"`
	Servers []*Server `yaml:"servers"`
	Striker *Striker  `yaml:"striker"`
}

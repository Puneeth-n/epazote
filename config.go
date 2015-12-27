package epazote

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Epazote struct {
	Config   Config
	Services map[string]Service
}

type Config struct {
	SMTP Email `yaml:"smtp"`
	HTTP Http  `yaml:"http"`
	Scan Scan  `yaml:"scan"`
}

type Email struct {
	Username string
	Password string
	Host     string
	Port     int
	Tls      bool
	Headers  map[string]string
}

type Http struct {
	Host     string
	pathPort int
}

type Scan struct {
	Paths                   []string
	Seconds, Minutes, Hours int
}

type Service struct {
	URL                     string
	Timeout                 int
	Seconds, Minutes, Hours int
	Log                     string
	Expect                  Expect
	IfStatus                map[string]Action `yaml:"if_status`
	IfHeader                map[string]Action `yaml:"if_header"`
}

type Expect struct {
	Status int
	Header map[string]string
	Body   string
	IfNot  Action `yaml:"if_not"`
}

type Action struct {
	Cmd    string
	Notify string
	Msg    string
}

func NewEpazote(file string) (*Epazote, error) {

	yml_file, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var ez Epazote

	if err := yaml.Unmarshal(yml_file, &ez); err != nil {
		return nil, err
	}

	return &ez, nil
}

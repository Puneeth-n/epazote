package config

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
}

type Email struct {
	Username string
	Password string
	Host     string
	Port     int
	Tls      bool
	Headers  map[string]string
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

func GetConfig(file string) (*Epazote, error) {

	yml_file, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var data Epazote

	if err := yaml.Unmarshal(yml_file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

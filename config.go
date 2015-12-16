package epazote

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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

func NewEpazote(file string) (*Epazote, error) {

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

func (ez *Epazote) CheckConfig() error {

	log.Printf("%# v", ez)
	return errors.New("path cannot be empty")

	return nil
}

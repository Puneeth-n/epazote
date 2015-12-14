package main

import (
	"fmt"
	"github.com/kr/pretty"
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

func main() {
	yml_file, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	var data Epazote

	if err := yaml.Unmarshal(yml_file, &data); err != nil {
		panic(err)
	}

	fmt.Printf("%# v", pretty.Formatter(data.Config.SMTP.Username))

	//	fmt.Println(data["service 1"].Expect.IfNot)

	//for k, v := range data["service 1"].Expect.Header {
	//	fmt.Println(k, v)
	//}
}

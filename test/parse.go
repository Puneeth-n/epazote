package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Services map[string]Service

type Service struct {
	URL                           string
	Timeout                       int
	Seconds, Minutes, Hours       int
	Log                           string
	If_not                        Expect
	If_status, If_header, If_body map[string]Action
}

type Action struct {
	Cmd    string
	Notify string
	Msg    string
}

type Expect struct {
	Status int
	Header map[string]string
	Body   string
	Action Action
}

func main() {
	yml_file, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	var data Services

	if err := yaml.Unmarshal(yml_file, &data); err != nil {
		panic(err)
	}

	fmt.Println(data["service 1"].If_not.Body)

	for k, v := range data["service 1"].If_not.Header {
		fmt.Println(k, v)
	}

}

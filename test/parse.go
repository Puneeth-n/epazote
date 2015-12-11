package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Services map[string]Service

type Service struct {
	Seconds   int
	Minutes   int
	Timeout   int
	URL       string
	Cmd       string
	Log       string
	Notify    string
	Msg       string
	If_header map[string]Action
	If_status map[string]Action
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

	var data Services

	if err := yaml.Unmarshal(yml_file, &data); err != nil {
		panic(err)
	}

	fmt.Println(data["service 1"].Seconds, data["service 1"].Cmd, data)

}

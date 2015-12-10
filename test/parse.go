package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Services map[string]interface{}

func main() {
	yml_file, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	var data Services

	if err := yaml.Unmarshal(yml_file, &data); err != nil {
		panic(err)
	}

	fmt.Println(data)

}

package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Get(s string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", s, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "epazote")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}

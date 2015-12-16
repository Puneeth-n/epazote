package epazote

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HTTPGet(s string) {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest("GET", s, nil)
	req.Header.Set("User-Agent", "epazote")

	resp, err := client.Do(req)
	if err != nil {
		//timeout check here
		log.Fatalln("timeout-----", err)
	}

	defer resp.Body.Close()
	//chunk := io.LimitReader(resp.Body, 0)
	chunk := io.LimitReader(resp.Body, 1<<20)
	body, err := ioutil.ReadAll(chunk)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body), len(body), resp)

}

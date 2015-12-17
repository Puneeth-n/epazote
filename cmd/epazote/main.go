package main

import (
	"flag"
	"fmt"
	ez "github.com/nbari/epazote"
	"log"
	"net/http"
	"os"
)

type HttpResponse struct {
	resp *http.Response
	err  error
}

const CRLF = "\r\n"

func main() {
	// f config file name
	var f = flag.String("f", "epazote.yml", "Epazote configuration file.")
	var c = flag.Bool("c", false, "Continue on errors.")

	flag.Parse()

	if _, err := os.Stat(*f); os.IsNotExist(err) {
		log.Fatalf("Cannot read file: %s, use -h for more info.\n\n", *f)
	}

	cfg, err := ez.NewEpazote(*f)
	if err != nil {
		log.Fatalln(err)
	}

	// # ----------------------------------------------------------------------------
	if *c {
	}

	ch := make(chan *HttpResponse, len(cfg.Services)) //buffered

	// check services before starting
	for k, v := range cfg.Services {
		go func(s ez.Service) {
			fmt.Printf("Checking URL for service: %s\n", k)
			resp, err := ez.Get(v)
			if err != nil {
				ch <- &HttpResponse{nil, err}
			}
			defer resp.Body.Close()
			ch <- &HttpResponse{resp, err}
		}(v)

		for {
			select {
			case r := <-ch:
				fmt.Printf("%s was fetched\n", r.resp)
			}
		}
	}

	//every := 60
	//if v.Seconds > 0 {
	//every = v.Seconds
	//} else if v.Minutes > 0 {
	//every = 60 * v.Minutes
	//} else if v.Hours > 0 {
	//every = 3600 * v.Hours
	//}
	//// test service
	//resp, err := ez.Get(v)
	//if err != nil {
	//log.Fatalf(ez.Red(fmt.Sprintf("Verify URL: %s for service: %s, error: %s", v.URL, k, err)))
	//}

	//log.Println(resp.StatusCode)

	//}

	//log.Printf("Supervising: %d services", len(cfg.Services))

}

// 	SendEmail(epazote.Config.SMTP)
//ez.Get("http://httpbin.org/get")

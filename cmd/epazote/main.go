package main

import (
	"flag"
	"github.com/kr/pretty"
	ez "github.com/nbari/epazote"
	"log"
	"net/url"
	"os"
)

func main() {

	// f config file name
	var f = flag.String("f", "epazote.yml", "Epazote configuration file.")
	var v = flag.Bool("v", false, "verbose mode")

	flag.Parse()

	if _, err := os.Stat(*f); os.IsNotExist(err) {
		log.Fatalf("Cannot read file: %s, use -h for more info.\n\n", *f)
	}

	c, err := ez.NewEpazote(*f)
	if err != nil {
		log.Fatalln(err)
	}

	if *v {
		log.Printf("%# v", pretty.Formatter(c))

		for k, v := range c.Services {
			log.Println("Service name: ", k)
			log.Println("URL:", v.URL)
			v.Every = 60
			if v.Seconds > 0 {
				v.Every = v.Seconds
			} else if v.Minutes > 0 {
				v.Every = 60 * v.Minutes
			} else if v.Hours > 0 {
				v.Every = 3600 * v.Hours
			}

			log.Printf("check every %d seconds", v.Every)

			// to a get to URL just to check if is recheable
			//			v.URL

		}
	}

	// 	SendEmail(epazote.Config.SMTP)
	//	ez.HTTPGet("http://httpbin.org/get")
}

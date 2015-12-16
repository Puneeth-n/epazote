package main

import (
	"flag"
	"github.com/kr/pretty"
	"log"
	"os"
)

func main() {

	// f config file name
	var f = flag.String("f", "epazote.yml", "Epazote configuration file.")
	var v = flag.Bool("v", false, "verbose, print configuration file.")

	flag.Parse()

	if _, err := os.Stat(*f); os.IsNotExist(err) {
		log.Fatalf("Cannot read file: %s, use -h for more info.\n\n", *f)
	}

	epazote, err := GetConfig(*f)
	if err != nil {
		panic(err)
	}

	if *v {
		log.Printf("%# v", pretty.Formatter(epazote))
	}

	//	fmt.Printf("%# v", epazote.Config.SMTP)

	// 	SendEmail(epazote.Config.SMTP)
	HTTPGet("http://httpbin.org/get")

}

package main

import (
	"flag"
	"github.com/kr/pretty"
	ez "github.com/nbari/epazote"
	"log"
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

	cfg, err := ez.NewEpazote(*f)
	if err != nil {
		log.Fatalln(err)
	}

	if *v {
		log.Printf("%# v", pretty.Formatter(cfg))
	}

	err = cfg.CheckConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//	fmt.Printf("%# v", epazote.Config.SMTP)

	// 	SendEmail(epazote.Config.SMTP)
	ez.HTTPGet("http://httpbin.org/get")

}

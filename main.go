package main

import (
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"github.com/nbari/epazote/config"
	"os"
)

func main() {

	// f config file name
	var f = flag.String("f", "epazote.yml", "Epazote configuration file.")
	var v = flag.Bool("v", false, "verbose, print configuration file.")

	flag.Parse()

	if _, err := os.Stat(*f); os.IsNotExist(err) {
		fmt.Printf("Can't read file: %s, use -h for more info.\n\n", *f)
		os.Exit(1)
	}

	epazote, err := config.GetConfig(*f)
	if err != nil {
		panic(err)
	}

	//	fmt.Printf("%#v", epazote.Config.SMTP.Username)

	if *v {
		fmt.Printf("%# v", pretty.Formatter(epazote))
		fmt.Println(epazote.Config.SMTP)
	}

}

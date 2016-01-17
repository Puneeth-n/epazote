package main

import (
	"flag"
	ez "github.com/nbari/epazote"
	"log"
	"os"
)

func main() {
	// f config file name
	var f = flag.String("f", "epazote.yml", "Epazote configuration file.")
	var c = flag.Bool("c", false, "Continue on errors.")

	flag.Parse()

	if _, err := os.Stat(*f); os.IsNotExist(err) {
		log.Fatalf("Cannot read file: %s, use -h for more info.\n\n", *f)
	}

	cfg, err := ez.New(*f)
	if err != nil {
		log.Fatalln(err)
	}

	if cfg == nil {
		log.Fatalln("Check config file sintax.")
	}

	// scan check config and clean paths
	err = cfg.CheckPaths()
	if err != nil {
		log.Fatalln(err)
	}

	// verify URL, we can't supervice unreachable services
	err = cfg.VerifyUrls()
	if err != nil {
		if !*c {
			log.Fatalln(err)
		}
		log.Println(err)
	}

	// check that at least a path or service are set
	err = cfg.PathsOrServices()
	if err != nil {
		log.Fatalln(err)
	}

	// verifyEMAIL recipients & headers
	err = cfg.VerifyEMAIL()
	if err != nil {
		log.Fatalln(err)
	}

	// create a Scheduler
	sk := ez.GetScheduler()

	cfg.Start(sk)
}

package main

import (
	"flag"
	ez "github.com/nbari/epazote"
	"log"
	"os"
)

const herb = "\U0001f33f"

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

	ch := ez.AsyncGet(cfg.Services)
	for i := 0; i < len(cfg.Services); i++ {
		x := <-ch
		if x.Err != nil {
			if !*c {
				log.Fatalf("%s - Verify URL: %q", ez.Red(x.Service), x.Err)
			}
			log.Printf("%s - Verify URL: %q", ez.Red(x.Service), x.Err)
		}
	}

	// create a new supervisor
	s := ez.NewSupervisor()

	// add services to supervisor
	for k, v := range cfg.Services {
		// how often to check for the service
		every := 60
		if v.Seconds > 0 {
			every = v.Seconds
		} else if v.Minutes > 0 {
			every = 60 * v.Minutes
		} else if v.Hours > 0 {
			every = 3600 * v.Hours
		}
		s.AddService(k, v, every)
	}

	log.Printf(ez.Green("Epazote %s   on %d services."), herb, len(cfg.Services))

	// block forever
	select {}
}

package main

import (
	"flag"
	//	"fmt"
	ez "github.com/nbari/epazote"
	"log"
	"os"
)

type ServiceStatus struct {
	err     error
	service string
}

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

	// check services before starting
	ch := make(chan ServiceStatus, len(cfg.Services))

	for k, v := range cfg.Services {
		go func(name string, s ez.Service) {
			resp, err := ez.Get(s)
			if err != nil {
				ch <- ServiceStatus{err, name}
				return
			}
			resp.Body.Close()
			ch <- ServiceStatus{nil, name}
		}(k, v)
	}

	for i := 0; i < len(cfg.Services); i++ {
		x := <-ch
		if x.err != nil {
			if !*c {
				log.Fatalf("%s - Verify URL: %q", ez.Red(x.service), x.err)
			}
			log.Printf("%s - Verify URL: %q", ez.Red(x.service), x.err)
		}
	}

	// add services to supervisor
	for k, v := range cfg.Services {
		every := 60
		if v.Seconds > 0 {
			every = v.Seconds
		} else if v.Minutes > 0 {
			every = 60 * v.Minutes
		} else if v.Hours > 0 {
			every = 3600 * v.Hours
		}
		ez.Supervice(k, v, every)
	}

	log.Printf(ez.Green("Epazote %s   supervising %d services."), herb, len(cfg.Services))

}

// 	SendEmail(epazote.Config.SMTP)

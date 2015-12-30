package main

import (
	"flag"
	ez "github.com/nbari/epazote"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	cfg, err := ez.New(*f)
	if err != nil {
		log.Fatalln(err)
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

	if len(cfg.Config.Scan.Paths) == 0 && len(cfg.Services) == 0 {
		log.Fatalln(ez.Red("No services to supervise or paths to scan."))
	}

	// create a Scheduler
	sk := ez.GetScheduler()

	// add services to supervisor
	for k, v := range cfg.Services {
		// how often to check for the service
		// default 1 minute
		every := 60
		if v.Seconds > 0 {
			every = v.Seconds
		} else if v.Minutes > 0 {
			every = 60 * v.Minutes
		} else if v.Hours > 0 {
			every = 3600 * v.Hours
		}

		// schedule service
		sk.AddScheduler(k, every, ez.Supervice(v))
	}

	if len(cfg.Config.Scan.Paths) > 0 {
		// default 5 minutes
		every := 300
		if cfg.Config.Scan.Seconds > 0 {
			every = cfg.Config.Scan.Seconds
		} else if cfg.Config.Scan.Minutes > 0 {
			every = 60 * cfg.Config.Scan.Minutes
		} else if cfg.Config.Scan.Hours > 0 {
			every = 3600 * cfg.Config.Scan.Hours
		}

		// set how often to scan
		s := new(ez.Scandir)
		for _, v := range cfg.Config.Scan.Paths {
			sk.AddScheduler(v, every, s.Scan(v))
		}
		log.Printf(ez.Green("Epazote %s   on %d services, scan paths: %s [pid: %d]."), herb, len(cfg.Services), strings.Join(cfg.Config.Scan.Paths, ","), os.Getpid())
	} else {
		log.Printf(ez.Green("Epazote %s   on %d services [pid: %d]."), herb, len(cfg.Services), os.Getpid())
	}

	// exit on signal
	block := make(chan os.Signal, 1)
	signal.Notify(block, os.Interrupt, os.Kill, syscall.SIGTERM)
	signalType := <-block
	signal.Stop(block)
	log.Printf("%q signal received.", signalType)
	sk.StopAll()
	log.Printf("Exiting.")
	os.Exit(0)
}

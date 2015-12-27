package main

import (
	"flag"
	ez "github.com/nbari/epazote"
	"github.com/nbari/epazote/scheduler"
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

	cfg, err := ez.NewEpazote(*f)
	if err != nil {
		log.Fatalln(err)
	}

	// scan check config
	if len(cfg.Config.Scan.Paths) > 0 {
		for _, d := range cfg.Config.Scan.Paths {
			if _, err := os.Stat(d); os.IsNotExist(err) {
				log.Fatalf("Verify that directory: %s, exists and is readable.", d)
			}
		}
	}

	// verify URL, we can't supervice unreachable services
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

	if len(cfg.Config.Scan.Paths) == 0 && len(cfg.Services) == 0 {
		log.Fatalln(ez.Red("No services to supervise or paths to scan."))
	}

	// create a Scheduler
	sk := scheduler.NewScheduler()

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
		for _, v := range cfg.Config.Scan.Paths {
			// how often to scan
			sk.AddScheduler(v, every, ez.Scandir(v))
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

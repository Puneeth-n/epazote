package main

import (
	"flag"
	"fmt"
	ez "github.com/nbari/epazote"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const herb = "\U0001f33f"

type Paths struct {
	paths []string
}

func (self *Paths) String() string {
	return fmt.Sprint(self.paths)
}

func (self *Paths) Set(value string) error {
	if len(self.paths) > 0 {
		return fmt.Errorf("The d flag is already set")
	}
	dirs := strings.Split(value, ",")
	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			return fmt.Errorf("Verify that directory: %s, exists and is readable.", d)
		}
		self.paths = append(self.paths, d)
	}
	return nil
}

func main() {

	var dFlag Paths

	// f config file name
	var f = flag.String("f", "epazote.yml", "Epazote configuration file.")
	var c = flag.Bool("c", false, "Continue on errors.")
	flag.Var(&dFlag, "d", "directory(s) to scan for epazote.yml, comma separated.")

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

	log.Printf(ez.Green("Epazote %s   on %d services [pid: %d]."), herb, len(cfg.Services), os.Getpid())

	// exit on signal
	block := make(chan os.Signal, 1)
	signal.Notify(block, os.Interrupt, os.Kill, syscall.SIGTERM)
	signalType := <-block
	signal.Stop(block)
	log.Printf("%q signal received.", signalType)
	s.Stop()
	log.Printf("Exiting.")
	os.Exit(0)
}

package main

import (
	"flag"
	"fmt"
	ez "github.com/nbari/epazote"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	fmt.Printf("%# v", cfg)

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

	// create a Scheduler
	sk := ez.GetScheduler()

	log.Printf("%s [pid: %d]", cfg.Start(sk), os.Getpid())

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

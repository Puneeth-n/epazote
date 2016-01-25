package epazote

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type IScheduler interface {
	AddScheduler(string, int, func())
	StopAll()
}

// Start Add services to scheduler
func (self *Epazote) Start(isk IScheduler, debug bool) {
	if debug {
		self.debug = true
	}

	for k, v := range self.Services {
		// Set service name
		v.Name = k

		// Status
		if v.Expect.Status < 1 {
			v.Expect.Status = 200
		}

		// rxBody
		if v.Expect.Body != "" {
			re := regexp.MustCompile(v.Expect.Body)
			v.Expect.body = re
		}

		if self.debug {
			if v.URL != "" {
				log.Printf(Green("Adding service: %s URL: %s"), v.Name, v.URL)
			} else {
				log.Printf(Green("Adding service: %s Test: %s"), v.Name, v.Test)
			}
		}

		// schedule the service
		isk.AddScheduler(k, GetInterval(60, v.Every), self.Supervice(v))
	}

	if len(self.Config.Scan.Paths) > 0 {
		for _, v := range self.Config.Scan.Paths {
			isk.AddScheduler(v, GetInterval(300, self.Config.Scan.Every), self.Scan(v))
		}
	}

	log.Printf("Epazote %s   on %d services, scan paths: %s [pid: %d]", herb, len(self.Services), strings.Join(self.Config.Scan.Paths, ","), os.Getpid())

	// stop until signal received
	start := time.Now()

	// loop forever
	block := make(chan os.Signal)

	signal.Notify(block, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		signalType := <-block
		switch signalType {
		case syscall.SIGUSR1, syscall.SIGUSR2:
			y, err := yaml.Marshal(&self)
			if err != nil {
				log.Printf("Error: %v", err)
			}
			l := `
    Gorutines: %d
    Alloc : %d
    Total Alloc: %d
    Sys: %d
    Lookups: %d
    Mallocs: %d
    Frees: %d
    Seconds in GC: %d
    Started on: %v
    Uptime: %v`

			runtime.NumGoroutine()
			s := new(runtime.MemStats)
			runtime.ReadMemStats(s)

			log.Printf("Config dump:\n%s---"+Green(l), y, runtime.NumGoroutine(), s.Alloc, s.TotalAlloc, s.Sys, s.Lookups, s.Mallocs, s.Frees, s.PauseTotalNs/1000000000, start.Format(time.RFC3339), time.Since(start))

		default:
			signal.Stop(block)
			log.Printf("%q signal received.", signalType)
			//			isk.StopAll()
			log.Println("Exiting.")
			os.Exit(0)
		}
	}
}

package epazote

import (
	"log"
	"os"
	"regexp"
	"strings"
)

type IScheduler interface {
	AddScheduler(string, int, func())
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

	log.Printf("Epazote %c   on %d services, scan paths: %s [pid: %d]", Icon(herb), len(self.Services), strings.Join(self.Config.Scan.Paths, ","), os.Getpid())
}

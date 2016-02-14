package epazote

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// Scan return func() to work with the scheduler
func (self *Epazote) Scan(dir string) func() {
	return func() {
		self.search(dir)
	}
}

// search walk through defined paths
func (self *Epazote) search(root string) error {
	find := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.Name() == "epazote.yml" {
			srv, err := ParseScan(path)
			if err != nil {
				return err
			}

			// get a Scheduler
			sk := GetScheduler()

			for k, v := range srv {
				if !IsURL(v.URL) {
					log.Printf("[%s] %s - Verify URL: %q", Red(path), k, v.URL)
					continue
				}

				// Set service name
				v.Name = k

				// Status
				if v.Expect.Status < 1 {
					v.Expect.Status = 200
				}

				// rxBody
				if v.Expect.Body != "" {
					re, err := regexp.Compile(v.Expect.Body)
					if err != nil {
						log.Printf("[%s] %s - Verify Body: %q - %q", Red(path), k, v.Expect.Body, err)
						continue
					}
					v.Expect.body = re
				}

				// Add/Update existing services
				if _, ok := self.Services[k]; !ok {
					self.Services[k] = v
				} else {
					old_status := self.Services[k].status
					print(old_status, "<------------")
					self.Services[k] = v
					self.Services[k].status = old_status
				}

				if self.debug {
					log.Printf(Green("Found epazote.yml in path: %s updating/adding service: %q"), path, k)
				}

				// schedule service
				sk.AddScheduler(k, GetInterval(60, v.Every), self.Supervice(self.Services[k]))
			}
		}
		return nil
	}

	// Walk over root using find func
	err := filepath.Walk(root, find)
	if err != nil {
		return err
	}

	return nil
}

package epazote

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Scandir struct{}

// Scan return func() to work with the scheduler
func (self Scandir) Scan(dir string) func() {
	return func() {
		self.search(dir)
	}
}

// search walk through defined paths
func (self Scandir) search(root string) error {
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

				// rxBody
				if body, ok := v.Expect.Body.(string); ok {
					re, err := regexp.Compile(body)
					if err != nil {
						log.Printf("[%s] %s - Verify Body: %q - %q", Red(path), k, body, err)
						continue
					}
					v.Expect.Body = *re
				}

				// schedule service
				sk.AddScheduler(k, GetInterval(60, v.Every), Supervice(v))
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

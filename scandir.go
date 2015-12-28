package epazote

import (
	"os"
	"path/filepath"
)

type (
	Services map[string]Service
	Scandir  struct{}
)

// Scan return func() to work with the scheduler
func (self Scandir) Scan(dir string) func() {
	return func() {
		self.search(dir)
	}
}

// search walk through defined paths
func (self Scandir) search(root string) {
	find := func(path string, f os.FileInfo, err error) error {
		if f.Name() == "epazote.yml" {
			return ParseScan(path)
		}
		return nil
	}

	filepath.Walk(root, find)
}

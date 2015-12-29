package epazote

import (
	"os"
	"path/filepath"
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
			return ParseScan(path)
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

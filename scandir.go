package epazote

import (
	"log"
	"os"
	"path/filepath"
)

type Scandir struct{}

func (self Scandir) Scan(dir string) func() {
	return func() {
		self.search(dir)
	}
}

func (self Scandir) search(root string) error {
	err := filepath.Walk(root, self.find)
	if err != nil {
		return err
	}
	return nil
}

func (self Scandir) find(path string, f os.FileInfo, err error) error {
	if f.Name() == "epazote.yml" {
		log.Printf("Visited: %s - %s\n", path, f.Name())
	}
	return nil
}

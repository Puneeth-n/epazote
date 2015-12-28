package epazote

import (
	"log"
	"os"
	"path/filepath"
)

type (
	Services map[string]Service
	Scandir  struct{}
)

func (self Scandir) Scan(dir string) func() {
	return func() {
		self.search(dir)
	}
}

func (self Scandir) search(root string) {
	err := filepath.Walk(root, self.find)
	if err != nil {
		log.Println(err)
	}
}

func (self Scandir) find(path string, f os.FileInfo, err error) error {
	if f.Name() == "epazote.yml" {
		return ParseScan(path)
	}
	return nil
}

package epazote

import (
	"fmt"
	"os"
	"path/filepath"
)

func find(path string, f os.FileInfo, err error) error {
	if f.Name() == "epazote.yml" {
		fmt.Printf("Visited: %s - %s\n", path, f.Name())
	}
	return nil
}

func Search(root string) error {
	err := filepath.Walk(root, find)
	return err
}

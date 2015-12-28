package epazote

import (
	//	"fmt"
	"testing"
)

func TestScan(t *testing.T) {
	s := new(Scandir)
	f := s.Scan("test")
	f()
}

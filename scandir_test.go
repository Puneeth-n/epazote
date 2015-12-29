package epazote

import (
	"reflect"
	"testing"
)

// just to cover
func TestScanReturnFunc(t *testing.T) {
	s := new(Scandir)
	f := s.Scan("test")
	ft := reflect.TypeOf(f).Kind()
	if ft != reflect.Func {
		t.Error("Expecting func()")
	} else {
		f()
	}
}

func TestScanSearchNonexistentRoot(t *testing.T) {
	s := new(Scandir)
	err := s.search("nonexistent")
	if err == nil {
		t.Error("Expecting: lstat nonexistent: no such file or directory")
	}
}

func TestScanSearch(t *testing.T) {
	s := new(Scandir)
	err := s.search("test")
	if err != nil {
		t.Error(err)
	}
}

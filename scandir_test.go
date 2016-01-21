package epazote

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

// just to cover
func TestScanReturnFunc(t *testing.T) {
	s := new(Epazote)
	f := s.Scan("test")
	ft := reflect.TypeOf(f).Kind()
	if ft != reflect.Func {
		t.Error("Expecting func()")
	} else {
		f()
	}
}

func TestScanSearchNonexistentRoot(t *testing.T) {
	s := new(Epazote)
	err := s.search("nonexistent")
	if err == nil {
		t.Error("Expecting: lstat nonexistent: no such file or directory")
	}
}

func TestScanSearch(t *testing.T) {
	s := new(Epazote)
	err := s.search("test")
	if err != nil {
		t.Error(err)
	}
}

func TestScanParseScanErr(t *testing.T) {
	dir := "./"
	prefix := "test-scan1-"

	d, err := ioutil.TempDir(dir, prefix)

	if err != nil {
		fmt.Println(err)
	}

	defer os.RemoveAll(d)

	f := []byte(`epazote
    - bad`)

	err = ioutil.WriteFile(fmt.Sprintf("%s/epazote.yml", d), f, 0644)

	s := new(Epazote)
	err = s.search(d)
	if err == nil {
		t.Error(err)
	}
}

func TestScanParseScanSearchOk(t *testing.T) {
	dir := "./"
	prefix := "test-scan2-"

	d, err := ioutil.TempDir(dir, prefix)

	if err != nil {
		fmt.Println(err)
	}

	defer os.RemoveAll(d)

	f := []byte(`
    service 1:
        url: http://about.epazote.io
        expect:
           body: "123"
`)

	err = ioutil.WriteFile(fmt.Sprintf("%s/epazote.yml", d), f, 0644)

	s := new(Epazote)
	err = s.search(d)
	if err != nil {
		t.Error(err)
	}
}

func TestScanParseScanSearchBadRegex(t *testing.T) {
	dir := "./"
	prefix := "test-scan2-"

	d, err := ioutil.TempDir(dir, prefix)

	if err != nil {
		fmt.Println(err)
	}

	defer os.RemoveAll(d)

	f := []byte(`
    service 1:
        url: http://about.epazote.io
        expect:
           body: ?(),
`)

	err = ioutil.WriteFile(fmt.Sprintf("%s/epazote.yml", d), f, 0644)

	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	s := new(Epazote)
	s.search(d)
	if err != nil {
		t.Error(err)
	}

	if buf.Len() == 0 {
		t.Error("Expecting log.Println error")
	}
	sk := GetScheduler()

	if len(sk.Schedulers) != 1 {
		t.Error("Expecting 1")
	}

}

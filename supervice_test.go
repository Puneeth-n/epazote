package epazote

import (
	"encoding/json"
	//	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestSuperviceStatusOk(t *testing.T) {
	check_s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-agent") != "epazote" {
			t.Error("Expecting User-agent: epazote")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer check_s.Close()
	log_s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-agent") != "epazote" {
			t.Error("Expecting User-agent: epazote")
		}
		decoder := json.NewDecoder(r.Body)
		var i map[string]interface{}
		err := decoder.Decode(&i)
		if err != nil {
			t.Error(err)
		}
		// check name
		if n, ok := i["Name"]; ok {
			if n != "s 1" {
				t.Errorf("Expecting  %q, got: %q", "s 1", n)
			}
		} else {
			t.Errorf("key not found: %q", "Name")
		}
		// check because
		if b, ok := i["Because"]; ok {
			if b != "Status: 200" {
				t.Errorf("Expecting: %q, got: %q", "Status: 200", b)
			}
		} else {
			t.Errorf("key not found: %q", "Because")
		}
		// check exit
		if e, ok := i["exit"]; ok {
			if e.(float64) != 0 {
				t.Errorf("Expecting: 0 got: %v", e.(float64))
			}
		} else {
			t.Errorf("key not found: %q", "exit")
		}
		wg.Done()
	}))
	defer log_s.Close()
	s := make(Services)
	s["s 1"] = Service{
		Name: "s 1",
		URL:  check_s.URL,
		Log:  log_s.URL,
		Expect: Expect{
			Status: 200,
		},
	}
	ez := &Epazote{
		Services: s,
	}
	wg.Add(1)
	ez.Supervice(s["s 1"])()
	wg.Wait()
}

package epazote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestSupervice(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-agent") != "epazote" {
			t.Error("Expecting User-agent: epazote")
		}
		fmt.Fprintln(w, "Hello, epazote")
	}))
	defer ts.Close()
	s := make(Services)
	s["s 1"] = Service{
		Name: "s 1",
		URL:  ts.URL,
		Expect: Expect{
			Status: 200,
		},
	}
	ez := &Epazote{
		Services: s,
	}
	fmt.Println(ez)
	f := ez.Supervice(s["s 1"])
	f()
}

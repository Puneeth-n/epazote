package epazote

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-agent") != "epazote" {
			t.Error("Expecting User-agent: epazote")
		}
		fmt.Fprintln(w, "Hello, epazote")
	}))
	defer ts.Close()

	res, err := HTTPGet(ts.URL, 3)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	if string(body) != "Hello, epazote\n" {
		t.Error("Expecting Hello, epazote")
	}

	if res.StatusCode != 200 {
		t.Error("Expecting StatusCode 200")
	}
}

func TestAsyngGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-agent") != "epazote" {
			t.Error("Expecting User-agent: epazote")
		}
		fmt.Fprintln(w, "Hello, epazote")
	}))
	defer ts.Close()
	s := make(Services)
	s["s 1"] = Service{
		URL: ts.URL,
	}
	ch := AsyncGet(s)
	for i := 0; i < len(s); i++ {
		x := <-ch
		if x.Err != nil {
			t.Error(x.Err)
		}
	}
}

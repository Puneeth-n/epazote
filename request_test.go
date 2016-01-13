package epazote

import (
	"testing"
)

func TestHTTPGet(t *testing.T) {
	_, err := HTTPGet("http://google.com", 3)
	if err != nil {
		t.Error(err)
	}
}

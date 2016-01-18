package epazote

import (
	"testing"
)

func TestVerifyBadEmail(t *testing.T) {
	cfg, err := New("test/epazote-bad-email.yml")
	if err != nil {
		t.Error(err)
	}
	err = cfg.VerifyEmail()
	if err == nil {
		t.Error("Expecting error")
	}
}

func TestVerifyEmail(t *testing.T) {
	cfg, err := New("test/epazote-email.yml")
	if err != nil {
		t.Error(err)
	}
	err = cfg.VerifyEmail()
	if err != nil {
		t.Errorf("Expecting error: %s", err)
	}
}

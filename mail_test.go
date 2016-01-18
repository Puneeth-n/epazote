package epazote

import (
	"gopkg.in/yaml.v2"
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

func TestVerifyEmailNoTo(t *testing.T) {
	cfg, err := New("test/epazote-email-noto.yml")
	if err != nil {
		t.Error(err)
	}
	err = cfg.VerifyEmail()
	if err == nil {
		t.Errorf("Expecting error: %s", err)
	}
}

func TestVerifyEmail2(t *testing.T) {
	var conf = `
config:
    smtp:
        username: username
        password: password
        server: smtp.server
        port: 587
        headers:
            from: from@email
            to: team@email
            subject: >
                [%s - %s], Service, Status
    services:
        service 1:
            expect:
                if_not:
                    notify: yes
                if_status:
                    502:
                        notify: yes
                if_header:
                    x-db-kaputt:
                        notify: yes
`
	var ez Epazote

	if err := yaml.Unmarshal([]byte(conf), &ez); err != nil {
		t.Error(err)
	}

}

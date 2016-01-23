package epazote

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"testing"
	"time"
)

// emailRecorder for testing
type emailRecorder struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
	msg  []byte
}

func mockSend(errToReturn error) (func(string, smtp.Auth, string, []string, []byte) error, *emailRecorder) {
	r := new(emailRecorder)
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		*r = emailRecorder{addr, a, from, to, msg}
		return errToReturn
	}, r
}

func TestEmail_SendSuccessful(t *testing.T) {
	c := &Email{}
	f, r := mockSend(nil)
	sender := &mailMan{c, f}
	body := "Hello World"
	err := sender.Send([]string{"me@example.com"}, []byte(body))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if string(r.msg) != body {
		t.Errorf("wrong message body.\n\nexpected: %\n got: %s", body, r.msg)
	}
}

func TestSendEmail(t *testing.T) {
	c := &Email{}
	f, r := mockSend(nil)
	sender := &mailMan{c, f}
	body := "Hello World"
	e := &Epazote{}
	e.SendEmail(sender, []string{"me@example.com"}, []byte(body))

	data, err := base64.StdEncoding.DecodeString(string(r.msg))
	if err != nil {
		t.Error(err)
	}
	if string(data) != body {
		fmt.Printf("%q\n", data)
	}
}

func TestReportNotify(t *testing.T) {
	c := Email{"username", "password", "server", 587, nil}
	f, r := mockSend(nil)
	sender := &mailMan{&c, f}
	ss := &Service{
		Name: "s 1",
		URL:  "http://about.epazote.io",
		Expect: Expect{
			Status: 200,
		},
	}
	a := &Action{Notify: "test@ejemplo.org"}
	e := &Epazote{}
	e.Config.SMTP = c

	e.Report(sender, ss, a, 1, 200, "because", "output")
	time.Sleep(10 * time.Millisecond)
	data, err := base64.StdEncoding.DecodeString(string(r.msg))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(data))
}

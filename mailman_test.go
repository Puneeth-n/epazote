package epazote

import (
	"encoding/base64"
	"fmt"
	"testing"
)

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
	//	err := sender.Send([]string{"me@example.com"}, []byte(body))
	err := e.SendEmail(sender, []string{"me@example.com"}, []byte(body))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(r.msg))
	if err != nil {
		t.Error(err)
	}
	if string(data) != body {
		fmt.Printf("%q\n", data)
	}
}

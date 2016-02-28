package epazote

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net/smtp"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// emailRecorder for testing
type emailRecorder struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
	msg  []byte
}

func mockSend(errToReturn error, wg *sync.WaitGroup) (func(string, smtp.Auth, string, []string, []byte) error, *emailRecorder) {
	r := new(emailRecorder)
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		defer wg.Done()
		*r = emailRecorder{addr, a, from, to, msg}
		return errToReturn
	}, r
}

func TestEmail_SendSuccessful(t *testing.T) {
	var wg sync.WaitGroup
	c := &Email{}
	f, r := mockSend(nil, &wg)
	sender := &mailMan{c, f}
	body := "Hello World"
	wg.Add(1)
	err := sender.Send([]string{"me@example.com"}, []byte(body))

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if string(r.msg) != body {
		t.Errorf("wrong message body.\n\nexpected: %s got: %s", body, r.msg)
	}
}

func TestSendEmail(t *testing.T) {
	var wg sync.WaitGroup
	c := &Email{}
	f, r := mockSend(nil, &wg)
	sender := &mailMan{c, f}
	body := "Hello World"
	e := &Epazote{}
	wg.Add(1)
	e.SendEmail(sender, []string{"me@example.com"}, "[name - exit]", []byte(body))

	data, err := base64.StdEncoding.DecodeString(string(r.msg))
	if err != nil {
		t.Error(err)
	}
	if string(data) != body {
		fmt.Printf("%q\n", data)
	}
}

func TestReportNotify(t *testing.T) {
	var wg sync.WaitGroup
	headers := map[string]string{
		"from": "epazote@domain.tld",
	}
	c := Email{"username", "password", "server", 587, headers, true}
	f, r := mockSend(nil, &wg)
	sender := &mailMan{&c, f}
	ss := &Service{
		Name: "s 1",
		URL:  "http://about.epazote.io",
		Expect: Expect{
			Status: 200,
		},
	}
	a := &Action{Notify: "test@ejemplo.org", Msg: []string{"OK", "NO OK"}, Emoji: []string{"1f621"}}
	e := &Epazote{}
	e.Config.SMTP = c

	wg.Add(1)
	e.Report(sender, ss, a, nil, 1, 200, "because", "output")
	wg.Wait()

	if r.addr != "server:587" {
		t.Errorf("Expecting %q got %q", "server:587", r.addr)
	}
	if r.from != "epazote@domain.tld" {
		t.Errorf("Expecting %q got %q", "epazote@domain.tld", r.from)
	}
	if r.to[0] != "test@ejemplo.org" {
		t.Errorf("Expecting %q got %q", "test@ejemplo.org", r.to[0])
	}

	crlf := []byte("\r\n\r\n")
	index := bytes.Index(r.msg, crlf)

	data := r.msg[index+len(crlf):]

	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		t.Error(err)
	}
}

func TestReportNotifyYes(t *testing.T) {
	var wg sync.WaitGroup
	buf.Reset()
	headers := map[string]string{
		"from":    "epazote@domain.tld",
		"to":      "test@ejemplo.org",
		"subject": "[name: name - exit - url - because]",
	}
	c := Email{"username", "password", "server", 587, headers, true}
	f, r := mockSend(errors.New("I love errors"), &wg)
	sender := &mailMan{&c, f}
	ss := &Service{
		Name: "s 1",
		URL:  "http://about.epazote.io",
		Expect: Expect{
			Status: 200,
		},
	}
	a := &Action{Notify: "yes", Msg: []string{"testing notifications"}}
	e := &Epazote{}
	e.Config.SMTP = c

	wg.Add(1)
	e.Report(sender, ss, a, nil, 1, 200, "because", "output")
	wg.Wait()

	if r.addr != "server:587" {
		t.Errorf("Expecting %q got %q", "server:587", r.addr)
	}
	if r.from != "epazote@domain.tld" {
		t.Errorf("Expecting %q got %q", "epazote@domain.tld", r.from)
	}
	if r.to[0] != "test@ejemplo.org" {
		t.Errorf("Expecting %q got %q", "test@ejemplo.org", r.to[0])
	}

	crlf := []byte("\r\n\r\n")
	index := bytes.Index(r.msg, crlf)

	data := r.msg[index+len(crlf):]

	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		t.Error(err)
	}

	if buf.Len() != 69 {
		t.Errorf("buf len not matching, got: %q", buf)
	}
}

func TestReportNotifySpam(t *testing.T) {
	var wg sync.WaitGroup
	buf.Reset()
	headers := map[string]string{
		"from":    "epazote@domain.tld",
		"to":      "test@ejemplo.org",
		"subject": "[name: name - exit - url - because]",
	}
	c := Email{"username", "password", "server", 587, headers, true}
	f, r := mockSend(errors.New("I love errors"), &wg)
	sender := &mailMan{&c, f}
	ss := &Service{
		Name: "s 1",
		URL:  "http://about.epazote.io",
		Expect: Expect{
			Status: 200,
		},
	}
	a := &Action{Notify: "yes", Msg: []string{"testing notifications"}}
	e := &Epazote{}
	e.Config.SMTP = c

	wg.Add(1)
	e.Report(sender, ss, a, nil, 1, 200, "because", "output")
	wg.Wait()

	if r.addr != "server:587" {
		t.Errorf("Expecting %q got %q", "server:587", r.addr)
	}
	if r.from != "epazote@domain.tld" {
		t.Errorf("Expecting %q got %q", "epazote@domain.tld", r.from)
	}
	if r.to[0] != "test@ejemplo.org" {
		t.Errorf("Expecting %q got %q", "test@ejemplo.org", r.to[0])
	}

	crlf := []byte("\r\n\r\n")
	index := bytes.Index(r.msg, crlf)

	data := r.msg[index+len(crlf):]

	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		t.Error(err)
	}

	if buf.Len() != 69 {
		t.Errorf("buf len not matching, got: %q", buf)
	}
}

func TestReportEmoji(t *testing.T) {
	var wg sync.WaitGroup
	buf.Reset()
	headers := map[string]string{
		"from":    "epazote@domain.tld",
		"to":      "test@ejemplo.org",
		"subject": "[name, because]",
	}
	c := Email{"username", "password", "server", 587, headers, true}
	f, r := mockSend(errors.New("I love errors"), &wg)
	sender := &mailMan{&c, f}
	ss := &Service{
		Name: "s 1",
		URL:  "http://about.epazote.io",
		Expect: Expect{
			Status: 200,
		},
	}
	a := &Action{Notify: "yes", Msg: []string{"testing notifications"}, Emoji: []string{"0"}}
	e := &Epazote{}
	e.Config.SMTP = c

	wg.Add(1)
	e.Report(sender, ss, a, nil, 1, 200, "because", "output")
	wg.Wait()

	if r.addr != "server:587" {
		t.Errorf("Expecting %q got %q", "server:587", r.addr)
	}
	if r.from != "epazote@domain.tld" {
		t.Errorf("Expecting %q got %q", "epazote@domain.tld", r.from)
	}
	if r.to[0] != "test@ejemplo.org" {
		t.Errorf("Expecting %q got %q", "test@ejemplo.org", r.to[0])
	}

	crlf := []byte("\r\n\r\n")
	index := bytes.Index(r.msg, crlf)

	data := r.msg[index+len(crlf):]

	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		t.Error(err)
	}

	if buf.Len() != 69 {
		t.Errorf("buf len not matching, got: %q", buf)
	}

	re := regexp.MustCompile(`Subject: \[s 1, because\]`)
	match := re.FindString(string(r.msg))
	if match == "" {
		t.Errorf("No subject found, expecting: %q", `Subject: [s 1, because]`)
	}
}

func TestReportCustomEmoji(t *testing.T) {
	var wg sync.WaitGroup
	buf.Reset()
	headers := map[string]string{
		"from":    "epazote@domain.tld",
		"to":      "test@ejemplo.org",
		"subject": "[name, because]",
	}
	c := Email{"username", "password", "server", 587, headers, true}
	f, r := mockSend(errors.New("I love errors"), &wg)
	sender := &mailMan{&c, f}
	ss := &Service{
		Name: "s 1",
		URL:  "http://about.epazote.io",
		Expect: Expect{
			Status: 200,
		},
	}
	a := &Action{Notify: "yes", Msg: []string{"testing notifications"}, Emoji: []string{"1F600", "1F621"}}
	e := &Epazote{}
	e.Config.SMTP = c

	wg.Add(1)
	e.Report(sender, ss, a, nil, 1, 200, "because", "output")
	wg.Wait()

	if r.addr != "server:587" {
		t.Errorf("Expecting %q got %q", "server:587", r.addr)
	}
	if r.from != "epazote@domain.tld" {
		t.Errorf("Expecting %q got %q", "epazote@domain.tld", r.from)
	}
	if r.to[0] != "test@ejemplo.org" {
		t.Errorf("Expecting %q got %q", "test@ejemplo.org", r.to[0])
	}

	crlf := []byte("\r\n\r\n")
	index := bytes.Index(r.msg, crlf)

	data := r.msg[index+len(crlf):]

	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		t.Error(err)
	}

	if buf.Len() != 69 {
		t.Errorf("buf len not matching, got: %q", buf)
	}

	re := regexp.MustCompile(`Subject.*`)
	match := re.FindString(string(r.msg))
	u := strings.Split(match, ": ")
	u[1] = strings.TrimSpace(u[1])
	dec := new(mime.WordDecoder)
	header, err := dec.DecodeHeader(u[1])
	if err != nil {
		t.Error(err)
	}
	x := fmt.Sprintf("%x", header)
	if x != "f09f98a120205b7320312c20626563617573655d" {
		t.Error(x)
	}
}

func TestReportMsg01(t *testing.T)   {}
func TestReportEmoji01(t *testing.T) {}

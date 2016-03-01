package epazote

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

// Log send log via HTTP POST to defined URL
func (self *Epazote) Log(s *Service, status []byte) {
	err := HTTPPost(s.Log, status, nil)
	if err != nil {
		log.Printf("Service %q - Error while posting to %q : %q", s.Name, s.Log, err)
	}
}

// Report create report to send via log/email
func (self *Epazote) Report(m MailMan, s *Service, a *Action, r *http.Response, e, status int, b, o string) {
	// set time
	t := time.Now().UTC().Format(time.RFC3339)

	// every (exit > 0) increment by one
	atomic.AddInt64(&s.status, 1)
	if e == 0 {
		s.status = 0
		if s.action != nil {
			a = s.action
		}
	}

	// create status report
	j, err := json.MarshalIndent(struct {
		*Service
		Exit    int    `json:"exit"`
		Status  int    `json:"status"`
		Output  string `json:"output,omitempty"`
		Because string `json:"because,omitempty"`
		When    string `json:"when"`
		Retries int    `json:"retries,omitempty"`
	}{s, e, status, o, b, t, s.retryCount}, "", "  ")

	if err != nil {
		log.Printf("Error creating report status for service %q: %s", s.Name, err)
		return
	}

	// debug
	if self.debug {
		// if available print the response headers
		var rHeader []string
		if r != nil {
			for k, _ := range r.Header {
				rHeader = append(rHeader, fmt.Sprintf("%s: %s", k, r.Header.Get(k)))
			}
			sort.Strings(rHeader)
		}
		if e == 0 {
			log.Printf(Green("Report: %s")+", Count: %d\n"+Yellow("Headers: \n%s\n"), j, s.status, strings.Join(rHeader, "\n"))
		} else {
			log.Printf(Red("Report: %s")+", Count: %d\n"+Yellow("Headers: \n%s\n"), j, s.status, strings.Join(rHeader, "\n"))

		}
	}

	if s.Log != "" {
		go self.Log(s, j)
	}

	// if no Action return
	if a == nil {
		return
	}

	// keys to be used in mail or in HTTP
	var parsed map[string]interface{}
	err = json.Unmarshal(j, &parsed)
	if err != nil {
		log.Printf("Error creating email report status for service %q: %s", s.Name, err)
		return
	}

	// sort the map
	var report_keys []string
	for k := range parsed {
		report_keys = append(report_keys, k)
	}
	sort.Strings(report_keys)

	// send email if action and only for the first error (avoid spam)
	if a.Notify != "" && s.status <= 1 {
		// store action on status so that when the service recovers
		// a notification can be sent to the previous recipients
		s.action = a

		if s.status == 0 {
			s.action = nil
		}

		// check if we can send emails
		if !self.Config.SMTP.enabled {
			log.Print(Red("Can't send email, no SMTP settings found."))
			return
		}

		// set To, recipients
		to := strings.Split(a.Notify, " ")
		if a.Notify == "yes" {
			to = strings.Split(self.Config.SMTP.Headers["to"], " ")
		}

		// prepare email body
		body := ""

		// based on the exit status select a  message to send
		// 0 - service OK
		// 1 - service failing
		msg := []string{"", ""}
		if len(a.Msg) > 1 {
			msg[0] = a.Msg[0]
			msg[1] = a.Msg[1]
		} else if len(a.Msg) == 1 {
			msg[0] = a.Msg[0]
		}

		body += fmt.Sprintf("%s %s%s", msg[s.status], CRLF, CRLF)

		// set subject _(because exit name output status url)_
		// replace the report status keys (json) in subject if present
		subject := self.Config.SMTP.Headers["subject"]
		for _, k := range report_keys {
			body += fmt.Sprintf("%s: %v %s", k, parsed[k], CRLF)
			subject = strings.Replace(subject, fmt.Sprintf("_%s_", k), fmt.Sprintf("%v", parsed[k]), 1)
		}

		// add emoji to subject
		emojis := []string{herb, shit}
		if len(a.Emoji) > 0 && a.Emoji[0] == "0" {
			emojis[0] = ""
			emojis[1] = ""
		} else if len(a.Emoji) == 1 {
			emojis[0] = a.Emoji[0]
		} else if len(a.Emoji) == 2 {
			emojis[0] = a.Emoji[0]
			emojis[1] = a.Emoji[1]
		}
		emoji := emojis[0]
		if s.status > 0 {
			emoji = emojis[1]
		}
		if emoji != "" {
			subject = mime.BEncoding.Encode("UTF-8", fmt.Sprintf("%c  %s", Icon(emoji), subject))
		}

		go self.SendEmail(m, to, subject, []byte(body))
	}

	// HTTP GET/POST based on exit status
	if len(a.HTTP) > 0 && s.status <= 1 {
		var h HTTP
		// if only one HTTP declared, use it if when service goes down (exit = 1)
		if len(a.HTTP) == 1 && s.status == 0 {
			return
		} else if len(a.HTTP) == 1 && s.status == 1 {
			h = a.HTTP[0]
		} else {
			h = a.HTTP[s.status]
		}
		if h.URL == "" {
			return
		}
		switch h.Method {
		case "post":
			// replace data with report_keys
			for _, k := range report_keys {
				h.Data = strings.Replace(h.Data, fmt.Sprintf("_%s_", k), fmt.Sprintf("%v", parsed[k]), 1)
			}
			go HTTPPost(h.URL, []byte(h.Data), h.Header)
		default:
			// replace url params with report_keys
			for _, k := range report_keys {
				h.URL = strings.Replace(h.URL, fmt.Sprintf("_%s_", k), fmt.Sprintf("%v", parsed[k]), 1)
			}
			go HTTPGet(h.URL, true, true, h.Header)
		}
		return
	}
}

package epazote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strings"
)

// Log send log via HTTP POST to defined URL
func (self *Epazote) Log(s *Service, status []byte) {
	err := HTTPPost(s.Log, status)
	if err != nil {
		log.Printf("Service %q - Error while posting to %q : %q", s.Name, s.Log, err)
	}
}

// Report create report to send via log/email
func (self *Epazote) Report(m MailMan, s *Service, a *Action, e, status int, b, o string) {
	// every exit 1 increment by one
	s.status++
	if e == 0 {
		s.status = 0
	}

	// create status report
	j, err := json.MarshalIndent(struct {
		*Service
		Exit    int    `json:"exit"`
		Status  int    `json:"status"`
		Output  string `json:"output,omitempty"`
		Because string `json:"because,omitempty"`
	}{s, e, status, o, b}, "", "  ")

	if err != nil {
		log.Printf("Error creating report status for service %q: %s", s.Name, err)
		return
	}

	// debug
	if self.debug {
		if e == 0 {
			log.Printf(Green("Report: %s"), j)
		} else {
			log.Printf(Red("Report: %s\nCount: %d"), j, s.status)
		}
	}

	if s.Log != "" {
		go self.Log(s, j)
	}

	// if no Action return
	if a == nil {
		return
	}

	// send email if action and only for the first error (avoid spam)
	if a.Notify != "" && s.status == 1 {
		var to []string
		if a.Notify == "yes" {
			to = strings.Split(self.Config.SMTP.Headers["to"], " ")
		} else {
			to = strings.Split(a.Notify, " ")
		}

		var parsed map[string]interface{}
		err := json.Unmarshal(j, &parsed)
		if err != nil {
			log.Printf("Error creating email report status for service %q: %s", s.Name, err)
			return
		}
		// sort the map
		var keys []string
		for k := range parsed {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// prepare email body
		body := ""
		if a.Msg != "" {
			body += fmt.Sprintf("%s %s%s", a.Msg, CRLF, CRLF)
		}

		// set subject
		subject := self.Config.SMTP.Headers["subject"]
		for _, k := range keys {
			body += fmt.Sprintf("%s: %v %s", k, parsed[k], CRLF)
			subject = strings.Replace(subject, k, fmt.Sprintf("%v", parsed[k]), 1)
		}

		go self.SendEmail(m, to, subject, []byte(body))
	}
}

// Do, execute the command in the if_not block
func (self *Epazote) Do(cmd *string, skip bool) string {
	if skip {
		return "Skipping cmd"
	}
	if len(*cmd) > 0 {
		args := strings.Fields(*cmd)
		out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return err.Error()
		}
		return string(out)
	}
	return "No defined cmd"
}

// Supervice check services
func (self *Epazote) Supervice(s Service) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Verify service %s options: %q", Red(s.Name), r)
			}
		}()

		// Mailman instance
		m := NewMailMan(&self.Config.SMTP)

		// skip do cmd, to avoid a loop
		var skip bool
		if s.status > s.Stop {
			skip = true
		}

		// Run Test if no URL
		// execute the Test cmd if exit > 0 execute the if_not cmd
		if s.URL == "" {
			args := strings.Fields(s.Test.Test)
			cmd := exec.Command(args[0], args[1:]...)
			var out bytes.Buffer
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				self.Report(m, &s, &s.Test.IfNot, 1, 0, fmt.Sprintf("Test cmd: %s", err), self.Do(&s.Test.IfNot.Cmd, skip))
			} else {
				self.Report(m, &s, nil, 0, 0, fmt.Sprintf("Test cmd: %s", out.String()), "")
			}
			return
		}

		// HTTP GET service URL
		res, err := HTTPGet(s.URL, s.Follow, s.Insecure, s.Timeout)
		if err != nil {
			self.Report(m, &s, &s.Expect.IfNot, 1, 0, fmt.Sprintf("GET: %s", err), self.Do(&s.Expect.IfNot.Cmd, skip))
			return
		}

		// Read Body first and close if not used
		if s.Expect.Body != "" {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Printf("Could not read Body for service %q: %s", Red(s.Name), err)
				return
			}
			r := s.Expect.body.FindString(string(body))
			if r == "" {
				self.Report(m, &s, &s.Expect.IfNot, 1, res.StatusCode, fmt.Sprintf("Body no regex match: %s", s.Expect.body.String()), self.Do(&s.Expect.IfNot.Cmd, skip))
			} else {
				self.Report(m, &s, nil, 0, res.StatusCode, fmt.Sprintf("Body regex match: %s", r), "")
			}
			return
		}

		// close body since will not be used anymore
		res.Body.Close()

		// if_status
		if s.IfStatus != nil {
			// chefk if there is an Action for the returned StatusCode
			if a, ok := s.IfStatus[res.StatusCode]; ok {
				self.Report(m, &s, &a, 1, res.StatusCode, fmt.Sprintf("Status: %d", res.StatusCode), self.Do(&a.Cmd, skip))
				return
			}
		}

		// if_header
		if s.IfHeader != nil {
			// return if true
			r := false
			for k, a := range s.IfHeader {
				if res.Header.Get(k) != "" {
					r = true
					self.Report(m, &s, &a, 1, res.StatusCode, fmt.Sprintf("Header: %s", k), self.Do(&a.Cmd, skip))
				}
			}
			if r {
				return
			}
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			self.Report(m, &s, &s.Expect.IfNot, 1, res.StatusCode, fmt.Sprintf("Status: %d", res.StatusCode), self.Do(&s.Expect.IfNot.Cmd, skip))
			return
		}

		// Header
		if s.Expect.Header != nil {
			for k, v := range s.Expect.Header {
				if !strings.HasPrefix(res.Header.Get(k), v) {
					self.Report(m, &s, &s.Expect.IfNot, 1, res.StatusCode, fmt.Sprintf("Header: %s", k), self.Do(&s.Expect.IfNot.Cmd, skip))
					return
				}
			}
		}

		// fin if all is ok
		if res.StatusCode == s.Expect.Status {
			self.Report(m, &s, nil, 0, res.StatusCode, fmt.Sprintf("Status: %d", res.StatusCode), "")
			return
		}
	}
}

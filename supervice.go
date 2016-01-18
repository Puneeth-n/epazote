package epazote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"os/exec"
	"regexp"
	"strings"
)

// Log exit(0|1) 0 successful, 1 failure
func (self *Epazote) Log(s *Service, exit int, because string, o ...string) (error, []byte) {
	if len(s.Log) > 0 {
		output := ""
		if len(o) > 0 {
			output = o[0]
		}

		// create json to send
		j, err := json.Marshal(struct {
			*Service
			Exit    int    `json:"exit"`
			Output  string `json:",omitempty"`
			Because string `json:",omitempty"`
		}{
			s,
			exit,
			output,
			because,
		})

		if err != nil {
			return err, nil
		}

		// If log
		err = HTTPPost(s.Log, j)
		if err != nil {
			log.Printf("Service %q - Error while posting to %q : %q", s.Name, s.Log, err)
		}
		return nil, j
	}
	return nil, nil
}

// Do, execute the command in the if_not block
func (self *Epazote) Do(s *Service, a *Action, because string) {
	cmd := a.Cmd
	j := []byte{}
	if len(cmd) > 0 {
		args := strings.Fields(cmd)
		out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			log.Printf("cmd error on service %q: %q", Red(s.Name), err)
		}
		err, j = self.Log(s, 1, because, string(out))
	}

	_, j = self.Log(s, 1, because)

	// send email
	if len(a.Notify) > 0 {
		var to []string
		for _, v := range strings.Split(a.Notify, " ") {
			e, err := mail.ParseAddress(v)
			if err != nil {
				continue
			}
			to = append(to, e.Address)
		}
		// send email
		log.Println(j)
		self.SendEmail(s, to, a.Msg)
	}
	return
}

// Supervice check services
func (self *Epazote) Supervice(s Service) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Verify service %q options: %q", Red(s.Name), r)
			}
		}()

		// Run Test if no URL
		// execute the Test cmd if exit > 0 execute the if_not cmd
		if s.URL == "" {
			args := strings.Fields(s.Test.Test)
			cmd := exec.Command(args[0], args[1:]...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				self.Do(&s, &s.Test.IfNot, fmt.Sprintf("Test cmd: %s", err))
				return
			}
			self.Log(&s, 0, fmt.Sprintf("Test cmd: %s", out.String()))
			return
		}

		// HTTP GET service URL
		res, err := HTTPGet(s.URL, s.Timeout)
		if err != nil {
			self.Do(&s, &s.Expect.IfNot, fmt.Sprintf("GET: %s", err))
			return
		}

		// Read Body first and close if not used
		if re, ok := s.Expect.Body.(regexp.Regexp); ok {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Printf("Could not read Body for service %q: %s", Red(s.Name), err)
				return
			}
			if re.FindString(string(body)) == "" {
				self.Do(&s, &s.Expect.IfNot, fmt.Sprintf("Body: %s", re.String()))
			}
			return
		}

		// close body since will not be used anymore
		res.Body.Close()

		// if_status
		if len(s.IfStatus) > 0 {
			if a, ok := s.IfStatus[res.StatusCode]; ok {
				self.Do(&s, &a, fmt.Sprintf("Status: %d", res.StatusCode))
			}
			return
		}

		// if_header
		if len(s.IfHeader) > 0 {
			for k, v := range s.IfHeader {
				if res.Header.Get(k) != "" {
					self.Do(&s, &v, fmt.Sprintf("Header: %q", k))
				}
			}
			return
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			self.Do(&s, &s.Expect.IfNot, fmt.Sprintf("Status: %d", res.StatusCode))
			return
		}

		// Header
		if len(s.Expect.Header) > 0 {
			for k, v := range s.Expect.Header {
				if res.Header.Get(k) != v {
					self.Do(&s, &s.Expect.IfNot, fmt.Sprint("Header: %s", k))
				}
				return
			}
		}

		// fin
		if res.StatusCode == s.Expect.Status {
			self.Log(&s, 0, fmt.Sprintf("Status: %d", res.StatusCode))
			return
		}

		log.Printf("Check service %q: %s", Red(s.Name), s.Expect.Status)
		return
	}
}

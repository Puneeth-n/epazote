package epazote

import (
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func TestCmd(c string) error {
	args := strings.Fields(c)
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (self *Epazote) Do(s *Service, a *Action) {
	cmd := a.Cmd
	if len(cmd) > 0 {
		args := strings.Fields(cmd)
		out, err := exec.Command(args[0], args[1:]...).Output()
		if err != nil {
			log.Printf("cmd error for service with URL: %s - %q:", Red(s.URL), err)
		}
		log.Printf("cmd output: %q", strings.TrimSpace(string(out)))
	}
	return
}

// Supervice check services
func (self *Epazote) Supervice(s Service) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Verify service options with URL: %s - %q", Red(s.URL), r)
			}
		}()

		// HTTP GET service URL
		res, err := HTTPGet(s.URL, s.Timeout)
		if err != nil {
			self.Do(&s, &s.Expect.IfNot)
			return
		}

		// Read Body first and close if not used
		if re, ok := s.Expect.Body.(regexp.Regexp); ok {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Printf("Could not read Body for service with URL: %s - %q:", Red(s.URL), err)
				return
			}
			if re.FindString(string(body)) == "" {
				self.Do(&s, &s.Expect.IfNot)
			}
			return
		}

		// close body since will not be used anymore
		res.Body.Close()

		// if_status
		if len(s.IfStatus) > 0 {
			if a, ok := s.IfStatus[res.StatusCode]; ok {
				self.Do(&s, &a)
			}
			return
		}

		// if_header
		if len(s.IfHeader) > 0 {
			for k, v := range s.IfHeader {
				if res.Header.Get(k) != "" {
					self.Do(&s, &v)
				}
			}
			return
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			self.Do(&s, &s.Expect.IfNot)
			return
		}

		// Header
		if len(s.Expect.Header) > 0 {
			for k, v := range s.Expect.Header {
				if res.Header.Get(k) != v {
					self.Do(&s, &s.Expect.IfNot)
				}
				return
			}
		}

		// fin
		if res.StatusCode == s.Expect.Status {
			log.Println("alles ok")
			return
		}

		log.Printf("Check service with URL: %s -%s", Red(s.URL), s.Expect.Status)
		return
	}
}

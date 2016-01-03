package epazote

import (
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func (s *Service) Do(a Action) {
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
func Supervice(s Service) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Verify service options with URL: %s - %q", Red(s.URL), r)
			}
		}()

		// HTTP GET service URL
		res, err := Get(s.URL, s.Timeout)
		if err != nil {
			s.Do(s.Expect.IfNot)
			return
		}

		// if_status
		if len(s.IfStatus) > 0 {
			if a, ok := s.IfStatus[res.StatusCode]; ok {
				s.Do(a)
			}
			return
		}

		// if_header
		if len(s.IfHeader) > 0 {
			for k, v := range s.IfHeader {
				if res.Header.Get(k) != "" {
					s.Do(v)
				}
			}
			return
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			s.Do(s.Expect.IfNot)
			return
		}

		// Header
		if len(s.Expect.Header) > 0 {
			for k, v := range s.Expect.Header {
				if res.Header.Get(k) != v {
					s.Do(s.Expect.IfNot)
				}
				return
			}
		}

		// Body
		if re, ok := s.Expect.Body.(regexp.Regexp); ok {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Printf("Could not read Body for service with URL: %s - %q:", Red(s.URL), err)
				return
			}
			if re.FindString(string(body)) == "" {
				s.Do(s.Expect.IfNot)
			}
			return
		}

		log.Printf("Check service with URL: %q", Red(s.URL))
		return
	}
}

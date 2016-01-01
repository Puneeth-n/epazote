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
		log.Println(string(out), err)
	}
	log.Println(cmd)
}

// Supervice check services
func Supervice(s Service) func() {
	return func() {
		// HTTP GET service URL
		res, err := Get(s.URL, s.Timeout)
		if err != nil {
			s.Do(s.Expect.IfNot)
		}

		defer res.Body.Close()

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
			}
			return
		}

		// Body
		if re, ok := s.Expect.Body.(regexp.Regexp); ok {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
			}
			if re.FindString(string(body)) == "" {
				s.Do(s.Expect.IfNot)
			}
			return
		}

		log.Printf("Check conf/body regex for service with url: %s", Red(s.URL))
		return
	}
}

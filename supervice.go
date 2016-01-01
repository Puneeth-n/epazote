package epazote

import (
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func Do(a Action) {
	cmd := a.Cmd
	if len(cmd) > 0 {
		args := strings.Fields(cmd)
		out, err := exec.Command(args[0], args[1:]...).Output()
		log.Println(out, err)
	}

}

// Supervice check services
func Supervice(s Service) func() {
	return func() {
		// HTTP GET service URL
		res, err := Get(s.URL, s.Timeout)
		if err != nil {
			Do(s.Expect.IfNot)
		}

		defer res.Body.Close()

		// if_status
		if len(s.IfStatus) > 0 {
			if a, ok := s.IfStatus[res.StatusCode]; ok {
				Do(a)
			}
			return
		}

		// if_header
		if len(s.IfHeader) > 0 {
			for k, v := range s.IfHeader {
				if res.Header.Get(k) != "" {
					Do(v)
				}
			}
			return
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			Do(s.Expect.IfNot)
			return
		}

		// Header
		if len(s.Expect.Header) > 0 {
			for k, v := range s.Expect.Header {
				if res.Header.Get(k) != v {
					Do(s.Expect.IfNot)
				}
			}
		}

		// Body
		if r, ok := s.Expect.Body.(*regexp.Regexp); ok {
			log.Printf("%# v", r)
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println(err, body)
			}
			if r.FindString(string(body)) == "" {
				Do(s.Expect.IfNot)
			}
		}

		log.Printf("%# v", s)
	}
}

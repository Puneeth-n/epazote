package epazote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func (self *Epazote) Log(s *Service, status []byte) {
	err := HTTPPost(s.Log, status)
	if err != nil {
		log.Printf("Service %q - Error while posting to %q : %q", s.Name, s.Log, err)
	}
}

func (self *Epazote) Report(s *Service, a *Action, e int, b string, o string) {

	// create status report
	j, err := json.Marshal(struct {
		*Service
		Exit    int    `json:"exit"`
		Output  string `json:"output,omitempty"`
		Because string `json:"because,omitempty"`
	}{s, e, o, b})

	if err != nil {
		log.Printf("Error creating report status for service %q: %s", s.Name, err)
		return
	}

	if s.Log != "" {
		go self.Log(s, j)
	}

	// if no Action return
	if a == nil {
		return
	}

	// action
	if a.Notify != "" {
		log.Println(j)
		//	go self.Notify(&s, s.Test.IfNot.Notify, status)
	}
}

// Do, execute the command in the if_not block
func (self *Epazote) Do(cmd *string) string {
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

func (self *Epazote) Notify(s *Service, to string, j []byte) {
}

// Supervice check services
func (self *Epazote) Supervice(s Service) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Verify service %s options: %q", Red(s.Name), r)
			}
		}()

		// Run Test if no URL
		// execute the Test cmd if exit > 0 execute the if_not cmd
		if s.URL == "" {
			args := strings.Fields(s.Test.Test)
			cmd := exec.Command(args[0], args[1:]...)
			var out bytes.Buffer
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				self.Report(&s, &s.Test.IfNot, 1, fmt.Sprintf("Test cmd: %s", err), self.Do(&s.Test.IfNot.Cmd))
			} else {
				self.Report(&s, nil, 0, fmt.Sprintf("Test cmd: %s", out.String()), "")
			}
			return
		}

		// HTTP GET service URL
		res, err := HTTPGet(s.URL, s.Timeout)
		if err != nil {
			self.Report(&s, &s.Expect.IfNot, 1, fmt.Sprintf("GET: %s", err), self.Do(&s.Expect.IfNot.Cmd))
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
			r := re.FindString(string(body))
			if r == "" {
				self.Report(&s, &s.Expect.IfNot, 1, fmt.Sprintf("Body no regex match: %s", re.String()), self.Do(&s.Expect.IfNot.Cmd))
			} else {
				self.Report(&s, nil, 0, fmt.Sprintf("Body regex match: %s", r), "")
			}
			return
		}

		// close body since will not be used anymore
		res.Body.Close()

		// if_status
		if len(s.IfStatus) > 0 {
			// chefk if there is an Action for the returned StatusCode
			if a, ok := s.IfStatus[res.StatusCode]; ok {
				self.Report(&s, &a, 1, fmt.Sprintf("Status: %d", res.StatusCode), self.Do(&a.Cmd))
				return
			}
		}

		// if_header
		if len(s.IfHeader) > 0 {
			//fmt.Printf("%#v", s.IfHeader)
			fmt.Println(len(s.IfHeader))
			for k, a := range s.IfHeader {
				fmt.Println("-oooooo", k, a)
				//	if res.Header.Get(k) != "" {
				//		self.Report(&s, &a, 1, fmt.Sprintf("Header: %s", k), self.Do(&a.Cmd))
				//	}
			}
			return
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			self.Report(&s, &s.Expect.IfNot, 1, fmt.Sprintf("Status: %d", res.StatusCode), self.Do(&s.Expect.IfNot.Cmd))
			return
		}

		// Header
		if len(s.Expect.Header) > 0 {
			for k, v := range s.Expect.Header {
				if res.Header.Get(k) != v {
					self.Report(&s, &s.Expect.IfNot, 1, fmt.Sprintf("Header: %s", k), self.Do(&s.Expect.IfNot.Cmd))
				}
				return
			}
		}

		// fin
		if res.StatusCode == s.Expect.Status {
			self.Report(&s, nil, 0, fmt.Sprintf("Status: %d", res.StatusCode), "")
			return
		}

		log.Printf("Check service %q: %s", Red(s.Name), s.Expect.Status)
		return
	}
}

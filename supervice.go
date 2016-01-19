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
	"sync"
)

// Status Service, exit(0|1) 0 successful, 1 failure, because (reason), output of command
func (self *Epazote) Status(s *Service, exit int, because string, o ...string) (error, []byte) {
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

	return nil, j
}

func (self *Epazote) Log(s *Service, status []byte, wg *sync.WaitGroup) {
	err := HTTPPost(s.Log, status)
	if err != nil {
		log.Printf("Service %q - Error while posting to %q : %q", s.Name, s.Log, err)
	}
	wg.Done()
}

// Do, execute the command in the if_not block
func (self *Epazote) Do(a *Action) string {
	cmd := a.Cmd
	if len(cmd) > 0 {
		args := strings.Fields(cmd)
		out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return err.Error()
		}
		return string(out)
	}
	return "No defined cmd"
}

func (self *Epazote) Notify(s *Service, to string, j []byte, wg *sync.WaitGroup) {
}

// Supervice check services
func (self *Epazote) Supervice(s Service) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Verify service %q options: %q", Red(s.Name), r)
			}
		}()

		// sync wait
		var wg sync.WaitGroup

		// Run Test if no URL
		// execute the Test cmd if exit > 0 execute the if_not cmd
		if s.URL == "" {
			args := strings.Fields(s.Test.Test)
			cmd := exec.Command(args[0], args[1:]...)
			var out bytes.Buffer
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				if err, status := self.Status(&s, 1, fmt.Sprintf("Test cmd: %s", err), self.Do(&s.Test.IfNot)); err != nil {
					log.Printf("Error creating status for service %q: %s", s.Name, err)
					return
				} else {
					// Log
					if s.Log != "" {
						wg.Add(1)
						go self.Log(&s, status, &wg)
					}
					// action
					if s.Test.IfNot.Notify != "" {
						wg.Add(1)
						go self.Notify(&s, s.Test.IfNot.Notify, status, &wg)
					}
					wg.Wait()
					return
				}
			}
			if err, status := self.Status(&s, 0, fmt.Sprintf("Test cmd: %s", out.String())); err != nil {
				log.Printf("Error creating status for service %q: %s", s.Name, err)
			} else {
				// Log
				if s.Log != "" {
					wg.Add(1)
					go self.Log(&s, status, &wg)
				}
				wg.Wait()
				return
			}
		}

		// HTTP GET service URL
		res, err := HTTPGet(s.URL, s.Timeout)
		if err != nil {
			err, status := self.Status(&s, 1, fmt.Sprintf("GET: %s", err), self.Do(&s.Expect.IfNot))
			if err != nil {
				log.Printf("Error creating status for service %q: %s", s.Name, err)
				return
			}
			// Log
			if s.Log != "" {
				wg.Add(1)
				go self.Log(&s, status, &wg)
			}
			// notify
			if s.Expect.IfNot.Notify != "" {
				wg.Add(1)
				go self.Notify(&s, s.Expect.IfNot.Notify, status, &wg)
			}
			wg.Wait()
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
				//				self.Do(&s, &s.Expect.IfNot, fmt.Sprintf("Body: %s", re.String()))
				err, status := self.Status(&s, 1, fmt.Sprintf("Body: %s", err), self.Do(&s.Expect.IfNot))
				if err != nil {
					log.Printf("Error creating status for service %q: %s", s.Name, err)
					return
				}
				// Log
				if s.Log != "" {
					wg.Add(1)
					go self.Log(&s, status, &wg)
				}
				// notify
				if s.Expect.IfNot.Notify != "" {
					wg.Add(1)
					go self.Notify(&s, s.Expect.IfNot.Notify, status, &wg)
				}
				wg.Wait()
				return
			}
			return
		}

		// close body since will not be used anymore
		res.Body.Close()

		// if_status
		if len(s.IfStatus) > 0 {
			if a, ok := s.IfStatus[res.StatusCode]; ok {
				err, status := self.Status(&s, 1, fmt.Sprintf("Status: %d", err), self.Do(&a))
				if err != nil {
					log.Printf("Error creating status for service %q: %s", s.Name, err)
					return
				}
				// Log
				if s.Log != "" {
					wg.Add(1)
					go self.Log(&s, status, &wg)
				}
				// notify
				if a.Notify != "" {
					wg.Add(1)
					go self.Notify(&s, a.Notify, status, &wg)
				}
				wg.Wait()
				return
			}
			return
		}

		// if_header
		if len(s.IfHeader) > 0 {
			for k, a := range s.IfHeader {
				if res.Header.Get(k) != "" {
					err, status := self.Status(&s, 1, fmt.Sprintf("Header: %q", err), self.Do(&a))
					if err != nil {
						log.Printf("Error creating status for service %q: %s", s.Name, err)
						return
					}
					// Log
					if s.Log != "" {
						wg.Add(1)
						go self.Log(&s, status, &wg)
					}
					// notify
					if a.Notify != "" {
						wg.Add(1)
						go self.Notify(&s, a.Notify, status, &wg)
					}
					wg.Wait()
					return
				}
			}
			return
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			err, status := self.Status(&s, 1, fmt.Sprintf("Status: %d", err), self.Do(&s.Expect.IfNot))
			if err != nil {
				log.Printf("Error creating status for service %q: %s", s.Name, err)
				return
			}
			// Log
			if s.Log != "" {
				wg.Add(1)
				go self.Log(&s, status, &wg)
			}
			// notify
			if s.Expect.IfNot.Notify != "" {
				wg.Add(1)
				go self.Notify(&s, s.Expect.IfNot.Notify, status, &wg)
			}
			wg.Wait()
			return
		}

		// Header
		if len(s.Expect.Header) > 0 {
			for k, v := range s.Expect.Header {
				if res.Header.Get(k) != v {
					err, status := self.Status(&s, 1, fmt.Sprintf("Header: %s", err), self.Do(&s.Expect.IfNot))
					if err != nil {
						log.Printf("Error creating status for service %q: %s", s.Name, err)
						return
					}
					// Log
					if s.Log != "" {
						wg.Add(1)
						go self.Log(&s, status, &wg)
					}
					// notify
					if s.Expect.IfNot.Notify != "" {
						wg.Add(1)
						go self.Notify(&s, s.Expect.IfNot.Notify, status, &wg)
					}
					wg.Wait()
				}
				return
			}
		}

		// fin
		if res.StatusCode == s.Expect.Status {
			err, status := self.Status(&s, 0, fmt.Sprintf("Status: %d", res.StatusCode))
			if err != nil {
				log.Printf("Error creating status for service %q: %s", s.Name, err)
			}
			// Log
			if s.Log != "" {
				wg.Add(1)
				go self.Log(&s, status, &wg)
			}
			wg.Wait()
			return
		}

		log.Printf("Check service %q: %s", Red(s.Name), s.Expect.Status)
		return
	}
}

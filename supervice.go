package epazote

import (
	"log"
	"os/exec"
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
		// to determine if body should be fetched or not
		get_body := false
		switch t := s.Expect.Body.(type) {
		case nil:
			log.Println(t)
		default:
			log.Println(t)
			get_body = true
		}

		// HTTP GET service URL
		res, err := Get(s.URL, s.Timeout)
		if err != nil {
			Do(s.Expect.IfNot)
		}

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
			}
			return
		}

		// Status
		if res.StatusCode != s.Expect.Status {
			Do(s.Expect.IfNot)
			return
		}

		// Header

		log.Println(get_body, res)
	}
}

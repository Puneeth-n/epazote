package epazote

import (
	"log"
	"os/exec"
	"strings"
)

// Exec
func Exec(cmd string) {
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

		res, err := Get(s.URL, s.Timeout)
		if err == nil {
			Exec(s.Expect.IfNot.Cmd)
		}
		log.Println(get_body, res)
	}
}

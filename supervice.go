package epazote

import (
	"log"
	//	"os/exec"
)

// Supervice check services
func Supervice(s Service) func() {
	return func() {
		get_body := false
		switch t := s.Expect.Body.(type) {
		case nil:
			log.Println(t)
		default:
			log.Println(t)
			get_body = true
		}

		log.Println(get_body)

	}
}

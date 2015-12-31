package epazote

import (
	"log"
	//	"os/exec"
)

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
		if err != nil {
			log.Println(err)
		}

		log.Println(get_body, res)

	}
}

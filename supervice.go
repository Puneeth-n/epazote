package epazote

import (
	"log"
)

func Supervice(service Service) func() {
	return func() {
		log.Println(service.Every, service.URL, service.Expect.Body)
	}
}

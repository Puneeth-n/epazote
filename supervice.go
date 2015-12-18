package epazote

import (
	"log"
	"time"
)

func Supervice(name string, s Service, every int) {
	e := time.Duration(every) * time.Second
	log.Println(name, s, e)
}

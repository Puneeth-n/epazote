package epazote

import (
	"log"
	"time"
)

type scheduler struct {
	t    <-chan time.Time
	quit chan struct{}
	f    func()
}

type Supervisor struct {
	services map[string]scheduler
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		services: make(map[string]scheduler),
	}
}

func (s *Supervisor) AddService(name string, service Service, every int) {
	e := time.Duration(every) * time.Second

	scheduler := scheduler{
		t:    time.NewTicker(e).C,
		quit: make(chan struct{}),
		f:    func() { s.Supervice(service) },
	}

	// add service
	s.services[name] = scheduler

	go func() {
		for {
			select {
			case <-scheduler.t:
				scheduler.f()
			case <-scheduler.quit:
				return
			}
		}
	}()
}

func (s *Supervisor) Supervice(service Service) {
	log.Println(service, service.URL)
}

func (s *Supervisor) Stop() {
	for k, v := range s.services {
		close(v.quit)
		log.Printf("Stoping: %s", k)
	}
}

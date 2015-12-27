package scheduler

import (
	"fmt"
	"log"
	"time"
)

type scheduler struct {
	name string
	t    <-chan time.Time
	quit chan struct{}
	f    func()
}

type Scheduler struct {
	schedulers map[string]scheduler
}

// NewScheduler returns a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		schedulers: make(map[string]scheduler),
	}
}

// AddScheduler calls a function every X seconds.
func (s *Scheduler) AddScheduler(name string, interval int, f func()) {
	e := time.Duration(interval) * time.Second

	scheduler := scheduler{
		name: name,
		t:    time.NewTicker(e).C,
		quit: make(chan struct{}),
		f:    f,
	}

	// add service
	s.schedulers[name] = scheduler

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

// Stop ends a specified scheduler.
func (s *Scheduler) Stop(name string) error {
	scheduler, ok := s.schedulers[name]

	if !ok {
		return fmt.Errorf("Scheduler: %s, does not exist.", name)
	}

	close(scheduler.quit)
	return nil
}

// StopAll ends all schedulers.
func (s *Scheduler) StopAll() {
	for k, v := range s.schedulers {
		close(v.quit)
		log.Printf("Stoping: %s", k)
	}
}

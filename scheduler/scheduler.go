package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	t    <-chan time.Time
	quit chan struct{}
	f    func()
}

type Schedulers struct {
	schedulers map[string]Scheduler
	sync.Mutex
}

// NewScheduler returns a new scheduler
func NewScheduler() *Schedulers {
	return &Schedulers{
		schedulers: make(map[string]Scheduler),
	}
}

// AddScheduler calls a function every X seconds.
func (s *Schedulers) AddScheduler(name string, interval int, f func()) {
	s.Lock()
	defer s.Unlock()

	e := time.Duration(interval) * time.Second

	scheduler := Scheduler{
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
func (s *Schedulers) Stop(name string) error {
	s.Lock()
	defer s.Unlock()

	scheduler, ok := s.schedulers[name]

	if !ok {
		return fmt.Errorf("Scheduler: %s, does not exist.", name)
	}

	close(scheduler.quit)
	return nil
}

// StopAll ends all schedulers.
func (s *Schedulers) StopAll() {
	s.Lock()
	defer s.Unlock()

	for k, v := range s.schedulers {
		close(v.quit)
		log.Printf("Stoping: %s", k)
	}
}

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
	Schedulers map[string]Scheduler
	sync.Mutex
}

// NewScheduler returns a new scheduler
func New() *Schedulers {
	return &Schedulers{
		Schedulers: make(map[string]Scheduler),
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

	// stop scheduler if exist
	if sk, ok := s.Schedulers[name]; ok {
		close(sk.quit)
	}

	// add service
	s.Schedulers[name] = scheduler

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

	sk, ok := s.Schedulers[name]

	if !ok {
		return fmt.Errorf("Scheduler: %s, does not exist.", name)
	}

	close(sk.quit)
	return nil
}

// StopAll ends all schedulers.
func (s *Schedulers) StopAll() {
	s.Lock()
	defer s.Unlock()

	for k, v := range s.Schedulers {
		close(v.quit)
		log.Printf("Stoping: %s", k)
	}
}

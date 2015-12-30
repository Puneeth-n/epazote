package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type scheduler struct {
	t    <-chan time.Time
	quit chan struct{}
	f    func()
}

type Scheduler struct {
	Schedulers map[string]scheduler
	sync.Mutex
}

// NewScheduler returns a new scheduler
func New() *Scheduler {
	return &Scheduler{
		Schedulers: make(map[string]scheduler),
	}
}

// AddScheduler calls a function every X seconds.
func (s *Scheduler) AddScheduler(name string, interval int, f func()) {
	s.Lock()
	defer s.Unlock()

	e := time.Duration(interval) * time.Second

	scheduler := scheduler{
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
func (s *Scheduler) Stop(name string) error {
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
func (s *Scheduler) StopAll() {
	s.Lock()
	defer s.Unlock()

	for k, v := range s.Schedulers {
		close(v.quit)
		log.Printf("Stoping: %s", k)
	}
}

package epazote

import (
	"fmt"
	"testing"
)

type fakeScheduler struct {
	name     string
	interval int
}

func (self *fakeScheduler) AddScheduler(name string, interval int, f func()) {
	self.name = name
	self.interval = interval
	fmt.Println(name, "<----")
}

func (self fakeScheduler) StopAll() {}

func TestStart(t *testing.T) {
	cfg, err := New("test/epazote-start.yml")
	if err != nil {
		t.Error(err)
	}
	err = cfg.PathsOrServices()
	if err != nil {
		t.Error(err)
	}
	sk := &fakeScheduler{}
	cfg.Start(sk, true)
}

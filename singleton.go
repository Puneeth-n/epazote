package epazote

import (
	sk "github.com/nbari/epazote/scheduler"
	"sync"
)

var instance *sk.Schedulers
var once sync.Once

func GetScheduler() *sk.Schedulers {
	once.Do(func() {
		instance = sk.New()
	})
	return instance
}

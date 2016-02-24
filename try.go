package epazote

import (
	"errors"
	"time"
)

// Func represents functions that can be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func Try(fn Func, interval int) error {
	var err error
	var cont bool
	sleep := time.Duration(interval) * time.Millisecond
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > 1<<16 {
			return errors.New("Exceeded retry limit")
		}
		time.Sleep(sleep)
	}
	return err
}

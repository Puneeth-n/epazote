package epazote

import (
	"github.com/kr/pretty"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func (self *Epazote) ProcessSignal() {
	// loop forever
	block := make(chan os.Signal)

	signal.Notify(block, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	sk := GetScheduler()

	for {
		signalType := <-block
		switch signalType {
		case syscall.SIGUSR1, syscall.SIGUSR2:
			log.Printf("%# v", pretty.Formatter(self))

			l := `
	Gorutines: %d",
	Alloc : %v
	Total Alloc: %v
	Lookups: %v
	Sys: %v
	Started on: %v
	Uptime: %v`

			runtime.NumGoroutine()
			s := new(runtime.MemStats)
			runtime.ReadMemStats(s)

			log.Printf(Green(l), runtime.NumGoroutine(), s.Alloc, s.TotalAlloc, s.Sys, s.Lookups, self.start.Format(time.RFC3339), time.Since(self.start))
		default:
			signal.Stop(block)
			log.Printf("%q signal received.", signalType)
			sk.StopAll()
			log.Println("Exiting.")
			os.Exit(0)
		}
	}
}

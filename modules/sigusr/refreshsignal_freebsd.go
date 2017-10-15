package sigusr

import (
	"os"
	"os/signal"
	"syscall"
)

func Handle(f func()) {
	sigsa := make(chan os.Signal, 1)
	signal.Notify(sigsa, syscall.SIGUSR1)
	go func() {
		for _ = range sigsa {
			go func() {
				f()
			}()
		}
	}()
}

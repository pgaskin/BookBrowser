package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func addRefreshSignalListener(server *Server) {
	sigsa := make(chan os.Signal, 1)
	signal.Notify(sigsa, syscall.SIGUSR1)
	go func() {
		for _ = range sigsa {
			go func() {
				log.Println("Booklist refresh triggered by SIGUSR1")
				server.RefreshBookIndex()
			}()
		}
	}()
}

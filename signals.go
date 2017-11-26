package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func waitSignals(done chan bool) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Println("Signal received: ", sig)
	done <- true // Time to say `Goodbye`
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func waitSignals(done chan bool) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Println("Signal received: ", sig)
	done <- true // Time to say `Goodbye`
}

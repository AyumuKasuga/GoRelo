package main

import (
	"os"
)

func main() {
	cmd := runProcess(os.Args[1:])
	defer gracefulShutdown(cmd)
	done := make(chan bool)
	go waitSignals(done)
	includeDirs := []string{"./"}
	go runWatch(includeDirs, cmd, os.Args[1:])
	<-done
}

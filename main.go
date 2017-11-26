package main

import (
	"os"
)

func main() {
	cProc := controlledProcess{
		command: os.Args[1:],
		cmd:     nil,
	}
	cProc.runProcess()
	defer cProc.gracefulShutdown()
	done := make(chan bool)
	go waitSignals(done)
	includeDirs := []string{"./"}
	go runWatch(includeDirs, cProc)
	<-done
}

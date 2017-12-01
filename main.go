package main

var cfg config
var cProc controlledProcess

func main() {
	cfg = config{}
	cfg.parseCmdArgs()

	cProc = controlledProcess{
		command: cfg.command,
	}
	cProc.runProcess()
	defer cProc.gracefulShutdown()
	done := make(chan bool)
	go waitSignals(done)
	go runWatch()
	<-done
}

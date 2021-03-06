package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

type controlledProcess struct {
	command []string
	cmd     *exec.Cmd
	exited  bool
}

func (p *controlledProcess) runProcess() {
	p.cmd = exec.Command(p.command[0], p.command[1:]...)
	p.cmd.Stdout, p.cmd.Stderr, p.cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	if err := p.cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go p.processWatcher()
}

func (p *controlledProcess) processWatcher() {
	err := p.cmd.Wait()
	if err != nil {
		log.Println(err)
	}
	log.Println("Program exited")
}

// Gently send SIGINT to the process, wait some time
// then if process still alive -- kill it.
func (p *controlledProcess) gracefulShutdown() {
	if p.cmd.ProcessState != nil {
		return // Process is already dead, nothing to do
	}
	log.Println("Sending SIGINT...")
	p.cmd.Process.Signal(os.Signal(cfg.killSignal))
	timePassed := time.Duration(0)
	for timePassed < cfg.killTimeout {
		if p.cmd.ProcessState != nil {
			return
		}
		<-time.After(time.Duration(time.Millisecond * 100))
		timePassed += time.Duration(time.Millisecond * 100)
	}
	if p.cmd.ProcessState == nil {
		log.Println("Process still alive, send SIGKILL")
		p.cmd.Process.Kill()
	}
}

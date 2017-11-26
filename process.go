package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func runProcess(command []string) (cmd *exec.Cmd) {
	cmd = exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go processWatcher(cmd)
	return
}

func gracefulShutdown(cmd *exec.Cmd) {
	if cmd.ProcessState != nil {
		return // Process is already dead, nothing to do
	}
	fmt.Println("Sending SIGINT...")
	cmd.Process.Signal(os.Signal(syscall.SIGINT))
	time.Sleep(time.Duration(time.Second))
	if cmd.ProcessState == nil {
		fmt.Println("Process still alive, send SIGKILL")
		cmd.Process.Kill()
	}
}

func processWatcher(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Program exited")
}

func waitSignals(sigs chan os.Signal, done chan bool) {
	sig := <-sigs
	fmt.Println("Signal received: ", sig)
	done <- true
}

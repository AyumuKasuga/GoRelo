package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	cmd := runProcess(os.Args[1:])
	sigs := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go waitSignals(sigs, done)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				fmt.Println("reloading...")
				gracefulShutdown(cmd)
				cmd = runProcess(os.Args[1:])
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}
	<-done
	gracefulShutdown(cmd)

}

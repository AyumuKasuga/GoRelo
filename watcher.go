package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func runWatch(includeDirs []string, cmd *exec.Cmd, cmdArgs []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range includeDirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}
			fmt.Println("reloading...")
			gracefulShutdown(cmd)
			cmd = runProcess(cmdArgs)
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func addRecursively(watcher *fsnotify.Watcher, dir string) {
	watcher.Add(dir)
	log.Println("DIR: ", dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		if file.IsDir() == true {
			addRecursively(watcher, filepath.Join(dir, file.Name()))
		}
	}
}

func runWatch(includeDirs []string, cmd *exec.Cmd, cmdArgs []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range includeDirs {
		addRecursively(watcher, dir)
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

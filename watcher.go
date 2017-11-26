package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func addRecursively(watcher *fsnotify.Watcher, dir string) {
	watcher.Add(dir)
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

func runWatch(includeDirs []string, cProc controlledProcess) {
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
			if event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Remove == fsnotify.Remove {
				watcher.Remove(event.Name)
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				addRecursively(watcher, event.Name)
			}

			fmt.Println("Something changes, reloading...")
			cProc.gracefulShutdown()
			cProc.runProcess()
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

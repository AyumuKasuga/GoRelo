package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

const reloadPeriod = time.Duration(time.Microsecond * 500)
const waitReloadPeriod = time.Duration(time.Microsecond * 200)

type reloader struct {
	lastEvent time.Time
}

func (r *reloader) periodicChecker() {
	for {
		<-time.After(reloadPeriod)
		if r.lastEvent.IsZero() == false && time.Since(r.lastEvent) > waitReloadPeriod {
			log.Println("Something changes, reloading...")
			cProc.gracefulShutdown()
			cProc.runProcess()
			r.lastEvent = time.Time{}
		}
	}
}

func (r *reloader) eventTrigger() {
	r.lastEvent = time.Now()
}

func addRecursively(watcher *fsnotify.Watcher, dir string) {
	for _, exclDir := range cfg.excludeFolders {
		if exclDir == dir {
			return
		}
	}
	watcher.Add(dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		// log.Println(err)
		// Hush little baby, user should not know about our little mistakes...
	}
	for _, file := range files {
		if file.IsDir() == true {
			addRecursively(watcher, filepath.Join(dir, file.Name()))
		}
	}
}

func runWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range cfg.watchFolders {
		addRecursively(watcher, dir)
	}

	rldr := reloader{}
	go rldr.periodicChecker()

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Remove == fsnotify.Remove {
				watcher.Remove(event.Name)
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				addRecursively(watcher, event.Name)
			}
			rldr.eventTrigger()
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

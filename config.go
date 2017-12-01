package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func getAbsDirs(dirs []string) (absDirs []string) {
	for _, dir := range dirs {
		absDir, _ := filepath.Abs(dir)
		absDirs = append(absDirs, absDir)
	}
	return
}

type config struct {
	watchFolders   []string
	excludeFolders []string
	killSignal     os.Signal
	killTimeout    time.Duration
	command        []string
}

func (c *config) parseCmdArgs() {
	flag.Usage = func() {
		fmt.Printf("Usage: gorelo [options] command \n\n")
		flag.PrintDefaults()
	}

	watchFolders := flag.String("w", "./", "Folders for watch (space separated)")
	excludeFolders := flag.String("e", "", "Exclude folders from watcher (space separated)")
	killTimeout := flag.Duration("t", time.Second, "Timeout before send SIGKILL")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	c.watchFolders = getAbsDirs(strings.Split(*watchFolders, " "))
	c.excludeFolders = getAbsDirs(strings.Split(*excludeFolders, " "))

	c.killTimeout = *killTimeout
	c.killSignal = syscall.SIGINT
	c.command = flag.Args()
	fmt.Println(c.killTimeout)
}

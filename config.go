package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type config struct {
	watchFolders   []string
	excludeFolders []string
	killSignal     os.Signal
	killTimeout    time.Duration
	command        []string
}

func (c *config) parseCmdArgs() {
	watchFolders := flag.String("w", "./", "Folders for watch (space separated)")
	excludeFolders := flag.String("e", "", "Do not watch these folders")
	killTimeout := flag.Duration("t", time.Second, "Timeout before send SIGKILL")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	c.watchFolders = strings.Split(*watchFolders, " ")
	c.excludeFolders = strings.Split(*excludeFolders, " ")
	c.killTimeout = *killTimeout
	c.killSignal = syscall.SIGINT
	c.command = flag.Args()
	fmt.Println(c.killTimeout)
}

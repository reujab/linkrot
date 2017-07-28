package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli"
)

var (
	waitgroup = new(sync.WaitGroup)
	semaphore = make(chan interface{}, 5)
)

var client = &http.Client{
	Timeout: time.Second * 10,
	CheckRedirect: func(*http.Request, []*http.Request) error {
		// don't redirect
		return http.ErrUseLastResponse
	},
}

func main() {
	app := cli.NewApp()
	app.Usage = "a command-line app that finds broken hyperlinks in files"
	app.Commands = []cli.Command{
		{
			Name:      "walk",
			Usage:     "walks a directory and checks every file",
			UsageText: app.Name + " walk [directory]",
			Action:    cmdWalk,
		},
	}
	app.Run(os.Args)
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

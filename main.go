package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli"
)

var (
	verbose    bool
	noHTTPS    bool
	noStatus   bool
	noRedirect bool
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
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "print successful requests",
			Destination: &verbose,
		},
		cli.BoolFlag{
			Name:        "no-https",
			Usage:       "don't check if domain supports https",
			Destination: &noHTTPS,
		},
		cli.BoolFlag{
			Name:        "no-status",
			Usage:       "don't check http status",
			Destination: &noStatus,
		},
		cli.BoolFlag{
			Name:        "no-redirect",
			Usage:       "don't check for redirects",
			Destination: &noRedirect,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "walk",
			Usage:     "Walks a directory and checks every file",
			UsageText: app.Name + " walk [directory]",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "exclude",
					Usage: "regular expressions to exclude",
				},
			},
			Action: cmdWalk,
		},
	}
	app.Run(os.Args)
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

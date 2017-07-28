package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"github.com/mvdan/xurls"
	"github.com/urfave/cli"
)

func cmdScan(ctx *cli.Context) {
	if len(ctx.Args()) == 0 {
		cli.ShowCommandHelpAndExit(ctx, "scan", -1)
	}

	for _, path := range ctx.Args() {
		file, err := os.Open(path)
		die(err)
		defer func() { die(file.Close()) }()

		var linesScanned int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			linesScanned++
			line := scanner.Text()

			// hyperlinks cannot be less than 11 characters
			// example of a short hyperlink: http://j.tl
			if len(line) < 11 {
				continue
			}

			for _, match := range xurls.Strict.FindAllString(line, -1) {
				uri, err := url.Parse(match)
				die(err)
				waitgroup.Add(1)
				go queue(fmt.Sprintf("%s:%d", path, linesScanned), uri)
			}
		}
	}

	waitgroup.Wait()
}

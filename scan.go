package main

import (
	"io/ioutil"
	"net/url"

	"github.com/mvdan/xurls"
	"github.com/urfave/cli"
)

func cmdScan(ctx *cli.Context) {
	if len(ctx.Args()) == 0 {
		cli.ShowCommandHelpAndExit(ctx, "scan", -1)
	}

	for _, arg := range ctx.Args() {
		buffer, err := ioutil.ReadFile(arg)
		die(err)
		for _, match := range xurls.Strict.FindAll(buffer, -1) {
			uri, err := url.Parse(string(match))
			die(err)
			waitgroup.Add(1)
			go queue(uri)
		}
	}

	waitgroup.Wait()
}

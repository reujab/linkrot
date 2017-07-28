package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/fatih/color"
)

var green = color.New(color.FgGreen)
var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)

func queue(file string, uri *url.URL) {
	semaphore <- nil
	defer waitgroup.Done()

	req, err := http.NewRequest("GET", uri.String(), nil)
	die(err)
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("%s %s %s\n", file, uri.String(), red.Sprint(err))
	} else if res.StatusCode >= 200 && res.StatusCode < 300 {
		if verbose {
			fmt.Printf("%s %s %s\n", file, uri.String(), green.Sprint(res.Status))
		}
	} else if res.StatusCode >= 300 && res.StatusCode < 400 {
		fmt.Printf("%s %s %s\n", file, uri.String(), yellow.Sprint(res.Status)+" -> "+res.Header.Get("Location"))
	} else {
		fmt.Printf("%s %s %s\n", file, uri.String(), red.Sprint(res.Status))
	}

	<-semaphore
}

package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	yellow = color.New(color.FgYellow)
	red    = color.New(color.FgRed)
)

var checkedURLs []string

func queue(file string, uri *url.URL) {
	defer waitgroup.Done()
	// return if url has already been checked
	for _, checkedURL := range checkedURLs {
		if uri.String() == checkedURL {
			return
		}
	}
	checkedURLs = append(checkedURLs, uri.String())

	semaphore <- nil

	req, err := http.NewRequest("HEAD", uri.String(), nil)
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

	if !noHTTPS && uri.Scheme != "https" {
		httpsURI := *uri
		httpsURI.Scheme = "https"
		req, err = http.NewRequest("HEAD", httpsURI.String(), nil)
		die(err)
		res, err = client.Do(req)

		if err == nil {
			fmt.Printf("%s %s %s\n", file, uri.String(), red.Sprint("use HTTPS"))
		}
	}

	<-semaphore
}

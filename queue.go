package main

import "net/url"

func queue(uri *url.URL) {
	semaphore <- nil
	defer waitgroup.Done()

	// TODO

	<-semaphore
}

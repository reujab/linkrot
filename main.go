package main

import (
	"net/http"
	"sync"
	"time"
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
	walk()
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

// Package webserver/demonstration is a example of how to use the Westeros
// webserver. Various features should be demonstrated here to help verify they
// work as intended and help others understand how to use the webserver package.
package main

import (
	"sync"

	"git.wreckerlabs.com/in/webserver"
)

const (
	port = "8080"
)

func main() {
	ws := webserver.New()

	// Example handler
	ws.GET("/", homeHandler)

	wg := &sync.WaitGroup{}

	// Webserver
	wg.Add(1)
	go func() {

		ws.Start(":8081")

		wg.Done()
	}()

	wg.Wait()
}

func homeHandler(event *webserver.Event) {
	event.HTML("demo", nil)
}

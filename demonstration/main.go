// Package webserver/demonstration is a example of how to use the Westeros
// webserver. Various features should be demonstrated here to help verify they
// work as intended and help others understand how to use the webserver package.
package main

import (
	"fmt"

	"git.wreckerlabs.com/in/webserver"
)

func main() {
	ws := webserver.New()

	fmt.Printf("%+v", ws)
}

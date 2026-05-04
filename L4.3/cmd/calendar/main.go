// Package main provides the entry point for the calendar application.
//
// It initializes the application by booting the App and starting the HTTP server.
package main

import (
	"L4.3/internal/app"
)

// main is the entry point of the application.
func main() {

	app.Boot().Run()

}

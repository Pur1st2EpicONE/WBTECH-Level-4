// Package main provides the entry point for the calendar application.
//
// It initializes the application by booting the App and starting the HTTP server.
package main

import (
	_ "L2.18/api/openapi-spec/docs"
	"L2.18/internal/app"
)

// main is the entry point of the application.
func main() {

	app.Boot().Run()

}

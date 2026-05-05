// Package main is the entry point of the GC and memory profiler service.
// It initializes the application and starts the HTTP server.
package main

import (
	"L4.4/internal/app"
)

// main starts the profiler application.
func main() {

	app.Boot().Run()

}

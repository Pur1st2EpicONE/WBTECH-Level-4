// Package main provides the entry point for the sort utility.
// It parses command-line flags and either runs the worker server or performs sorting.
package main

import (
	"fmt"

	"L4.2/internal/flags"
	"L4.2/internal/server"
	"L4.2/internal/sorter"
)

// main parses flags and either starts the worker server or runs the sorting operation.
func main() {

	flags, err := flags.Parse()
	if err != nil {
		sorter.LogFatal(fmt.Errorf("failed to parse flags: %w", err), "")
	}

	if flags.Serve {
		fmt.Println("Starting worker on port", flags.Port)
		if err := server.ListenAndServe(flags.Port); err != nil {
			sorter.LogFatal(fmt.Errorf("server run failed: %w", err), "")
		}
		return
	}

	sorter.Sort(flags)

}

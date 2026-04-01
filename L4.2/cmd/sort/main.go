// Package main provides the entry point for the custom sort utility.
// It delegates all sorting logic, flag parsing, and file handling to the
// internal sorter package, keeping the main function minimal.
package main

import (
	"L4.2/internal/sorter"
)

func main() {

	sorter.Sort()

}

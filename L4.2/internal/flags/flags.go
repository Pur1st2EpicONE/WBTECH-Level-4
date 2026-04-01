// Package flags provides parsing and storage of command-line flags
// for the sort utility, similar to GNU sort options.
package flags

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/spf13/pflag"
)

// Flags holds all available command-line options for the sorting utility.
type Flags struct {
	K          bool          // Sort via a key (-k)
	N          bool          // Numeric sort (-n)
	R          bool          // Reverse order sort (-r)
	U          bool          // Unique output (-u)
	M          bool          // Month sort (-M)
	B          bool          // Ignore leading blanks (-b)
	C          bool          // Check if input is sorted without sorting (-c)
	H          bool          // Human-readable numeric sort (-h)
	kRaw       string        // Raw -k option value, used to determine the key column
	Blanks     *atomic.Int32 // Counter of blank lines, relevant when -k is active
	ClmnToSort int           // Index of the column to sort when -k is used
}

// Parse parses the command-line flags and returns a Flags struct.
// Returns an error if parsing or -k key processing fails.
func Parse() (*Flags, error) {

	flags := &Flags{Blanks: new(atomic.Int32)}

	pflag.StringVarP(&flags.kRaw, "key", "k", "", "sort via a key; KEYDEF gives location and type")
	pflag.BoolVarP(&flags.N, "numeric-sort", "n", false, "compare according to string numerical value")
	pflag.BoolVarP(&flags.R, "reverse", "r", false, "reverse the result of comparisons")
	pflag.BoolVarP(&flags.U, "unique", "u", false, "with -c, check for strict ordering;\n  without -c, output only the first of an equal run")
	pflag.BoolVarP(&flags.M, "month-sort", "M", false, "compare (unknown) < 'JAN' < ... < 'DEC'")
	pflag.BoolVarP(&flags.B, "ignore-leading-blanks", "b", false, "ignore leading blanks")
	pflag.BoolVarP(&flags.C, "check", "c", false, "check for sorted input; do not sort")
	pflag.BoolVarP(&flags.H, "human-numeric-sort", "h", false, "compare human readable numbers (e.g., 2K 1G)")

	pflag.Parse()

	if err := kCheck(&flags.K, flags.kRaw, &flags.ClmnToSort); err != nil {
		return nil, fmt.Errorf("-k: %w", err)
	}

	return flags, nil

}

// kCheck validates and parses the -k key option.
// It sets the flagK boolean to true if -k is provided, and extracts the column to sort.
// Returns an error if the column number cannot be converted to an integer.
func kCheck(flagK *bool, kString string, clmnToSort *int) error {
	if kString != "" {
		*flagK = true
		fields := strings.FieldsSeq(kString)
		for number := range fields {
			if kInt, err := strconv.Atoi(number); err != nil {
				return fmt.Errorf("unable to convert number using %w", err)
			} else {
				if kInt > 0 {
					*clmnToSort = kInt - 1
				} else {
					*clmnToSort = 0
				}
				break
			}
		}
	}
	return nil
}

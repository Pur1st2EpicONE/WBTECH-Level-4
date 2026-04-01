// Package comparator provides line comparison functions for the sort utility,
// supporting numeric, human-readable, month-based, and column/key-based sorting.
package comparator

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"L4.2/internal/flags"
)

// Compare returns a comparison function for two strings based on the provided Flags.
func Compare(flags *flags.Flags) func(a, b string) int {

	return func(line1, line2 string) int {

		if flags.B {
			line1 = strings.TrimLeft(line1, " \t")
			line2 = strings.TrimLeft(line2, " \t")
		}

		if flags.K {

			fields1 := strings.Fields(line1)
			fields2 := strings.Fields(line2)
			var column1, column2 string

			if len(fields1) > flags.ClmnToSort {
				column1 = fields1[flags.ClmnToSort]
			}

			if len(fields2) > flags.ClmnToSort {
				column2 = fields2[flags.ClmnToSort]
			}

			return CompareLines(column1, column2, flags)

		}

		return CompareLines(line1, line2, flags)

	}

}

// CompareLines compares two lines according to the given Flags
func CompareLines(line1, line2 string, flags *flags.Flags) int {

	if flags.B {
		line1 = strings.TrimLeft(line1, " \t")
		line2 = strings.TrimLeft(line2, " \t")
	}

	var cmpRes int

	switch {
	case flags.N:
		cmpRes = numericComparison(line1, line2)
	case flags.H:
		cmpRes = readableComparison(line1, line2)
	case flags.M:
		cmpRes = monthsComparison(line1, line2)
	default:
		cmpRes = stringComparison(line1, line2)
	}

	if flags.R {
		return -cmpRes
	}
	return cmpRes

}

// numericComparison compares lines based on their numeric value at the start.
// If a line does not start with a number, it falls back to lexicographic comparison.
func numericComparison(line1, line2 string) int {

	number1, isNumber1 := getFirstInt(line1)
	number2, isNumber2 := getFirstInt(line2)

	switch {
	case isNumber1 && isNumber2:
		return compareNumbers(number1, number2, line1, line2)
	case isNumber1:
		return 1
	case isNumber2:
		return -1
	default:
		return cmp.Compare(line1, line2)
	}

}

// firstIntRegexp matches the first integer at the beginning of a string.
var firstIntRegexp = regexp.MustCompile(`^\s*([0-9]+)`)

// getFirstInt extracts the first integer at the start of a line.
// Returns the number and a boolean indicating whether a number was found.
func getFirstInt(line string) (int64, bool) {

	firstInt := firstIntRegexp.FindStringSubmatch(line)
	if len(firstInt) < 2 {
		return 0, false
	}

	number, err := strconv.ParseInt(firstInt[1], 10, 64)
	return number, err == nil

}

// compareNumbers compares two numbers and falls back to line comparison if equal.
func compareNumbers(number1, number2 int64, line1, line2 string) int {

	if res := cmp.Compare(number1, number2); res != 0 {
		return res
	}
	return cmp.Compare(line1, line2)

}

// readableComparison compares human-readable numbers.
func readableComparison(line1 string, line2 string) int {

	number1, isNumber1 := convertToInt(line1)
	number2, isNumber2 := convertToInt(line2)

	switch {
	case isNumber1 == nil && isNumber2 == nil:
		if number1 < number2 {
			return -1
		} else if number1 > number2 {
			return 1
		}
		return stringComparison(line1, line2)
	case isNumber1 == nil && isNumber2 != nil:
		return 1
	case isNumber1 != nil && isNumber2 == nil:
		return -1
	default:
		return stringComparison(line1, line2)
	}

}

// IntRe matches a human-readable number at the start of a string, possibly
// with a decimal part and an optional unit suffix (K, M, G, T, P, E).
var IntRe = regexp.MustCompile(`^\s*([0-9]+(?:\.[0-9]+)?)([KkMmGgTtPpEe])?`)

// convertToInt parses human-readable numbers and returns their value as int64.
func convertToInt(line string) (int64, error) {

	number := IntRe.FindStringSubmatch(line)
	if number == nil {
		return 0, fmt.Errorf("no number found in the string")
	}

	value, err := strconv.ParseFloat(number[1], 64)
	if err != nil {
		return 0, err
	}

	return int64(value * multiplier(number)), nil

}

// multiplier returns the factor corresponding to the unit suffix.
func multiplier(number []string) float64 {

	switch strings.ToLower(number[2]) {
	case "k":
		return 1024
	case "m":
		return 1024 * 1024
	case "g":
		return 1024 * 1024 * 1024
	case "t":
		return 1024 * 1024 * 1024 * 1024
	case "p":
		return 1024 * 1024 * 1024 * 1024 * 1024
	case "e":
		return 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	}

	return 1

}

// monthsComparison compares lines assuming month names.
func monthsComparison(line1 string, line2 string) int {

	months := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
		"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}

	month1 := monthToInt(months, line1)
	month2 := monthToInt(months, line2)

	if month1 == 0 && month2 == 0 {
		return stringComparison(line1, line2)
	}

	if month1 == 0 {
		return -1
	}

	if month2 == 0 {
		return 1
	}

	if month1 != month2 {
		if month1 < month2 {
			return -1
		}
		return 1
	}

	return stringComparison(line1, line2)

}

// monthToInt converts a month abbreviation to an integer.
func monthToInt(months map[string]int, line string) int {

	line = strings.TrimLeft(line, " \t")

	if len(line) < 3 {
		return 0
	}

	if month, found := months[strings.ToLower(line[:3])]; found {
		return month
	}

	return 0

}

// stringComparison compares two strings lexicographically.
func stringComparison(line1 string, line2 string) int {

	if line1 < line2 {
		return -1
	} else if line1 > line2 {
		return 1
	}
	return 0

}

// CheckSorted verifies if the input from scanner is sorted according to Flags.
// Returns false if a disorder is found, printing a message similar to GNU sort.
func CheckSorted(scanner *bufio.Scanner, fileName string, flags *flags.Flags) (sorted bool, err error) {

	if fileName != "" {

		file, openErr := os.Open(fileName)
		if openErr != nil {
			return false, openErr
		}
		scanner = bufio.NewScanner(file)

		defer func() {
			if closeErr := file.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

	} else {
		fileName = "-"
	}

	var line string
	var prev string
	var lineNum int

	for scanner.Scan() {

		line = scanner.Text()
		lineNum++

		if lineNum == 1 {
			prev = line
			continue
		}

		if flags.U && line == prev {
			fmt.Fprintf(os.Stderr, "sort: %s:%d: disorder: %s\n", fileName, lineNum, line)
			return false, nil
		}

		if CompareLines(line, prev, flags) < 0 {
			fmt.Fprintf(os.Stderr, "sort: %s:%d: disorder: %s\n", fileName, lineNum, line)
			return false, nil
		}

		prev = line

	}

	return true, scanner.Err()

}

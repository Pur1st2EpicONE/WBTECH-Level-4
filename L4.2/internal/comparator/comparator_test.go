package comparator

import (
	"bufio"
	"strings"
	"testing"

	"L4.2/internal/flags"
)

func TestStringComparison(t *testing.T) {

	tests := []struct {
		line1, line2 string
		expected     int
	}{
		{"qwe", "qwe", 0},
		{"abc", "abd", -1},
		{"abd", "abc", 1},
	}

	for _, tt := range tests {
		if result := stringComparison(tt.line1, tt.line2); result != tt.expected {
			t.Errorf("stringComparison(%q, %q) = %d, expected = %d", tt.line1, tt.line2, result, tt.expected)
		}
	}

}

func TestNumericComparison(t *testing.T) {

	tests := []struct {
		line1, line2 string
		expected     int
	}{
		{"2", "1", 1},
		{"1", "2", -1},
		{"aboba", "1", -1},
		{"1", "amogus", 1},
		{"abc", "def", -1},
	}

	for _, tt := range tests {
		if result := numericComparison(tt.line1, tt.line2); result != tt.expected {
			t.Errorf("numericComparison(%q, %q) = %d, expected = %d", tt.line1, tt.line2, result, tt.expected)
		}
	}

}

func TestReadableComparison(t *testing.T) {

	tests := []struct {
		line1, line2 string
		expected     int
	}{
		{"2K", "1K", 1},
		{"100M", "100M", 0},
		{"1G", "500M", 1},
		{"qwe", "1K", -1},
		{"1K", "qwe", 1},
	}

	for _, tt := range tests {
		if result := readableComparison(tt.line1, tt.line2); result != tt.expected {
			t.Errorf("readableComparison(%q, %q) = %d, expected = %d", tt.line1, tt.line2, result, tt.expected)
		}
	}

}

func TestMonthsComparison(t *testing.T) {

	tests := []struct {
		line1, line2 string
		expected     int
	}{
		{"Jan", "Feb", -1},
		{"jan", "feb", -1},
		{"May", "May", 0},
		{"Dec", "nov", 1},
		{"qwe", "Jan", -1},
		{"Jan", "QWE", 1},
	}

	for _, tt := range tests {
		if result := monthsComparison(tt.line1, tt.line2); result != tt.expected {
			t.Errorf("monthsComparison(%q, %q) = %d, expected = %d", tt.line1, tt.line2, result, tt.expected)
		}
	}

}
func TestCheckSorted(t *testing.T) {

	flags := &flags.Flags{}

	t.Run("sorted", func(t *testing.T) {
		data := "a\nb\nc\n"
		scanner := bufio.NewScanner(strings.NewReader(data))

		sorted, err := CheckSorted(scanner, "", flags)
		if err != nil {
			t.Errorf("CheckSorted returned unexpected error: %v", err)
		}
		if !sorted {
			t.Errorf("CheckSorted(%q) = false, expected true", data)
		}
	})

	t.Run("unsorted", func(t *testing.T) {
		data := "a\nc\nb\n"
		scanner := bufio.NewScanner(strings.NewReader(data))

		sorted, err := CheckSorted(scanner, "", flags)
		if err != nil {
			t.Errorf("CheckSorted returned unexpected error: %v", err)
		}
		if sorted {
			t.Errorf("CheckSorted(%q) = true, expected false", data)
		}
	})

}

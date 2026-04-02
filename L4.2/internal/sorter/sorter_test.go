package sorter

import (
	"bufio"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"L4.2/internal/config"
	"L4.2/internal/flags"
)

func TestSortChunkAndSaveChunk(t *testing.T) {

	lines := []string{"b", "a", "c"}
	expectedLines := []string{"a", "b", "c"}

	SortChunk(lines, new(flags.Flags))

	for i := range lines {
		if lines[i] != expectedLines[i] {
			t.Errorf("sortChunk: lines[%d] = %q, expected = %q", i, lines[i], expectedLines[i])
		}
	}

	fileName, err := saveChunk(lines)
	if err != nil {
		t.Fatalf("saveChunk returned error: %v", err)
	}

	defer func() {
		if rmErr := os.Remove(fileName); rmErr != nil {
			t.Logf("failed to remove temp file: %v", rmErr)
		}
	}()

	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("failed to read saved chunk: %v", err)
	}

	expected := "a\nb\nc\n"
	if string(result) != expected {
		t.Errorf("saveChunk: result = %q, expected = %q", string(result), expected)
	}

}

func TestSplitToChunks(t *testing.T) {

	data := &Data{
		Flags:  new(flags.Flags),
		Config: &config.Config{ChunkSize: 2, Workers: 2},
	}

	lines := []string{"q", "w", "e"}
	scanner := bufio.NewScanner(strings.NewReader(strings.Join(lines, "\n")))

	if err := splitToChunks(data, scanner); err != nil {
		t.Fatalf("splitToChunks returned error: %v", err)
	}

	if len(data.Chunks) == 0 {
		t.Errorf("splitToChunks: no chunks created, expected > 0")
	}

	for _, chunkFile := range data.Chunks {
		if rmErr := os.Remove(chunkFile); rmErr != nil {
			t.Logf("failed to remove chunk file: %v", rmErr)
		}
	}

}

func TestGetLesserLineIdx(t *testing.T) {

	flags := &flags.Flags{}
	lines := []string{"b", "a", "c"}

	idx := getLesserLineIdx(-1, 1, lines, flags)
	if idx != 1 {
		t.Errorf("getLesserLineIdx: result = %d, expected = %d", idx, 1)
	}

}

func TestUpdateLines(t *testing.T) {

	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	defer func() {
		if rmErr := os.Remove(file.Name()); rmErr != nil {
			t.Logf("failed to remove temp file: %v", rmErr)
		}
	}()

	if _, err := file.WriteString("q\nw\n"); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek temp file: %v", err)
	}

	scanner := bufio.NewScanner(file)
	scanners := []*bufio.Scanner{scanner}
	lines := []string{""}
	alive := []bool{true}

	updateLines(scanners, lines, alive, 0)
	if lines[0] != "q" || !alive[0] {
		t.Errorf("updateLines: lines[0] = %q, alive[0] = %v, expected = %q, %v", lines[0], alive[0], "q", true)
	}

}

func TestPrintBlanks(t *testing.T) {
	amount := &atomic.Int32{}
	amount.Store(3)
	printBlanks(amount)
}

func TestMergePrep(t *testing.T) {

	result := []string{"a\n", "b\n"}
	chunkFiles := make([]string, len(result))

	for i, line := range result {

		f, err := os.CreateTemp("", "chunk_test")
		if err != nil {
			t.Fatalf("failed to create temp chunk: %v", err)
		}

		if _, err := f.WriteString(line); err != nil {
			t.Fatalf("failed to write chunk: %v", err)
		}

		if err := f.Close(); err != nil {
			t.Fatalf("failed to close chunk file: %v", err)
		}

		chunkFiles[i] = f.Name()

		defer func(name string) {
			if rmErr := os.Remove(name); rmErr != nil {
				t.Logf("failed to remove chunk file: %v", rmErr)
			}
		}(f.Name())

	}

	scnrs := make([]*bufio.Scanner, len(chunkFiles))
	lines := make([]string, len(chunkFiles))
	alive := make([]bool, len(chunkFiles))
	files := make([]*os.File, len(chunkFiles))

	if err := mergePrep(scnrs, files, lines, alive, chunkFiles); err != nil {
		t.Fatalf("mergePrep returned error: %v", err)
	}

	for i, l := range lines {
		if l == "" {
			t.Errorf("mergePrep: lines[%d] = %q, expected non-empty string", i, l)
		}
	}

}

func TestCleanup(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "cleanup_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	chunks := []string{tmpFile.Name()}
	if err := cleanup([]*os.File{}, chunks); err != nil {
		t.Errorf("cleanup returned error: %v", err)
	}

	if _, err := os.Stat(tmpFile.Name()); !os.IsNotExist(err) {
		t.Errorf("cleanup: file %q was not deleted", tmpFile.Name())
	}

}

func TestPrepareLine(t *testing.T) {

	flags := &flags.Flags{K: true, ClmnToSort: 1}
	line := "Why are you reading this"

	result := prepareLine(line, flags)
	expected := "are"

	if result != expected {
		t.Errorf("prepareLine(%q) = %q, expected = %q", line, result, expected)
	}

}

func TestOutput(t *testing.T) {

	var prevLine string
	var printedEmpty bool
	flags := &flags.Flags{U: true}

	output("test", &prevLine, flags, &printedEmpty)
	if prevLine != "test" {
		t.Errorf("output: prevLine = %q, expected = %q", prevLine, "test")
	}

}

func TestDeleteChunks(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "file_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	chunks := []string{tmpFile.Name()}
	if err := deleteChunks(chunks); err != nil {
		t.Errorf("deleteChunks returned error: %v", err)
	}

	if _, err := os.Stat(tmpFile.Name()); !os.IsNotExist(err) {
		t.Errorf("deleteChunks: file %q was not removed", tmpFile.Name())
	}

}

// Package sorter provides functions for sorting text data, either from files
// or standard input, supporting chunked, in-memory, and optional distributed sorting.
package sorter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"L4.2/internal/comparator"
	"L4.2/internal/config"
	"L4.2/internal/master"

	"L4.2/internal/flags"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

const (
	inputNotSorted   = 1
	internalError    = 2
	defaultChunkSize = 1000
)

// Data represents the complete context of a sorting operation,
// including configuration, command-line flags, file metadata,
// and intermediate sorting chunks.
type Data struct {
	Chunks   []string       // Paths to temporary chunk files
	Flags    *flags.Flags   // Parsed command-line flags
	FileName string         // Current input file name, used for concise error message formatting
	Config   *config.Config // Runtime configuration settings
	Sorted   bool           // Indicates whether the input data is already sorted; used with the -c flag.
}

// Sort is the main entry point for the sorting process.
// It parses input, processes chunks, performs sorting, merging, and cleanup.
func Sort(flags *flags.Flags) {

	data, err := processInput(flags)
	if err != nil {
		LogFatal(err, data.FileName)
	}

	if !data.Flags.C {

		scnrs := make([]*bufio.Scanner, len(data.Chunks))
		files := make([]*os.File, len(data.Chunks))
		lines := make([]string, len(data.Chunks))
		alive := make([]bool, len(data.Chunks))

		if err := mergePrep(scnrs, files, lines, alive, data.Chunks); err != nil {
			LogFatal(err, data.FileName)
		}

		if err := mergeSort(scnrs, lines, alive, data.Flags); err != nil {
			LogFatal(err, data.FileName)
		}

		if err := cleanup(files, data.Chunks); err != nil {
			LogFatal(err, data.FileName)
		}

	}

}

// processInput determines the input source and initializes runtime configuration.
func processInput(f *flags.Flags) (*Data, error) {

	data := &Data{Sorted: true, Flags: f}

	cfg, err := config.Load()
	if err != nil || cfg.ChunkSize < 1 {
		data.Config = &config.Config{ChunkSize: defaultChunkSize}
	} else {
		data.Config = cfg
	}

	files := pflag.Args()
	if len(files) > 0 {
		return data, processFiles(data, files)
	}
	return data, processStdIn(data)

}

// processFiles processes one or more input files sequentially.
func processFiles(data *Data, files []string) error {

	for i := range files {

		data.FileName = files[i]

		file, err := os.Open(data.FileName)
		if err != nil {
			return err // intentionally not wrapped: the error is handled by LogFatal, which formats it in the style of GNU sort using errors.Is checks.
		}

		if err = processData(data, bufio.NewScanner(file)); err != nil {
			return err // intentionally not wrapped
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("unable to process %s file: %w", files[i], err)
		}

	}

	return nil

}

// processStdIn processes standard input when no files are provided.
func processStdIn(data *Data) error {

	if err := processData(data, bufio.NewScanner(os.Stdin)); err != nil {
		return err // intentionally not wrapped: the error is handled by LogFatal, which formats it in the style of GNU sort using errors.Is checks.
	}
	return nil

}

// processData sorts or checks sortedness of the provided scanner data
// depending on flags, and splits input into chunks.
func processData(data *Data, scanner *bufio.Scanner) error {

	if data.Flags.C {

		var err error
		data.Sorted, err = comparator.CheckSorted(scanner, data.FileName, data.Flags)
		if err != nil {
			return err // intentionally not wrapped: the error is handled by LogFatal, which formats it in the style of GNU sort using errors.Is checks.
		}

		if !data.Sorted {
			os.Exit(inputNotSorted) // Exit with code 1 to mimic GNU sort behavior when input is not sorted (-c flag).
		}

	}

	return splitToChunks(data, scanner)

}

// splitToChunks divides input into chunks and distributes work among workers.
func splitToChunks(data *Data, scanner *bufio.Scanner) error {

	workers := runtime.GOMAXPROCS(data.Config.Workers)
	chunksQueue := make(chan []string)

	var mu sync.Mutex
	var g errgroup.Group

	for range workers {
		g.Go(func() error {
			return processChunks(data, chunksQueue, &mu)
		})
	}

	lines := 0
	chunkBuilder := make([]string, 0, data.Config.ChunkSize)

	for scanner.Scan() {
		chunkBuilder = append(chunkBuilder, scanner.Text())
		lines++
		if lines == data.Config.ChunkSize {
			chunk := make([]string, len(chunkBuilder))
			copy(chunk, chunkBuilder)
			chunksQueue <- chunk
			chunkBuilder = chunkBuilder[:0] // Reset chunkBuilder to reuse it for the next chunk without reallocating
			lines = 0
		}
	}

	if len(chunkBuilder) > 0 {
		chunksQueue <- chunkBuilder // Send any remaining lines in chunkBuilder as the final chunk
	}

	close(chunksQueue)

	if err := g.Wait(); err != nil {
		return fmt.Errorf("worker malfunction: %w", err)
	}

	if err := scanner.Err(); err != nil {
		return err // intentionally not wrapped: the error is handled by LogFatal, which formats it in the style of GNU sort using errors.Is checks.
	}

	return nil

}

// processChunks sorts individual chunks and writes them to temporary files.
func processChunks(data *Data, chunksQueue <-chan []string, mu *sync.Mutex) error {

	for chunk := range chunksQueue {

		if data.Flags.Nodes != "" {

			nodes := strings.Split(data.Flags.Nodes, ",")

			sortedChunk, err := master.RemoteSort(chunk, nodes, data.Flags, data.Flags.Quorum)
			if err != nil {
				return err
			}

			filename, err := saveChunk(sortedChunk)
			if err != nil {
				return err
			}

			mu.Lock()
			data.Chunks = append(data.Chunks, filename)
			mu.Unlock()

			continue

		}

		if err := localSort(chunk, data, mu); err != nil {
			return err
		}

	}

	return nil

}

// localSort sorts a chunk and writes it to a temporary file.
func localSort(chunk []string, data *Data, mu *sync.Mutex) error {

	SortChunk(chunk, data.Flags)

	filename, err := saveChunk(chunk)
	if err != nil {
		return err
	}

	mu.Lock()
	data.Chunks = append(data.Chunks, filename)
	mu.Unlock()

	return nil

}

// SortChunk sorts a slice of lines in-place according to flags.
func SortChunk(lines []string, flags *flags.Flags) {

	slices.SortStableFunc(lines, comparator.Compare(flags))

	if flags.K {
		for _, line := range lines {
			if line == "" {
				flags.Blanks.Add(1)
				continue
			}
		}
	}

}

// saveChunk writes sorted lines into a temporary file.
func saveChunk(chunk []string) (fileName string, err error) {

	tempFile, err := os.CreateTemp(".", "chunk_")
	if err != nil {
		return "", fmt.Errorf("unable to create temporary file: %w", err)
	}

	defer func() {
		if closeErr := tempFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	writer := bufio.NewWriter(tempFile)
	for _, line := range chunk {
		if _, writeErr := writer.WriteString(line + "\n"); writeErr != nil {
			return "", fmt.Errorf("unable to write string to temporary file: %w", writeErr)
		}
	}

	if flushErr := writer.Flush(); flushErr != nil {
		return "", fmt.Errorf("unable to flush buffered data: %w", flushErr)
	}

	return tempFile.Name(), nil

}

// mergePrep prepares scanners and files for the merging stage.
func mergePrep(scnrs []*bufio.Scanner, files []*os.File, lines []string, alive []bool, chunks []string) error {

	for i, file := range chunks {

		currentFile, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("unable to open temporary file %s: %w", file, err)
		}

		files[i] = currentFile
		scnrs[i] = bufio.NewScanner(currentFile)

		if scnrs[i].Scan() {
			lines[i] = scnrs[i].Text()
			alive[i] = true
		}

	}

	return nil

}

// mergeSort merges sorted chunks into final sorted output.
func mergeSort(scnrs []*bufio.Scanner, lines []string, alive []bool, flags *flags.Flags) error {

	if flags.K && !flags.R {
		printBlanks(flags.Blanks)
	}

	var prevLine string
	var printedEmpty bool
	var lesserLineIdx int

	for {

		lesserLineIdx = -1

		for currentIdx := range lines {
			if !alive[currentIdx] {
				continue
			}
			lesserLineIdx = getLesserLineIdx(lesserLineIdx, currentIdx, lines, flags)
		}

		if lesserLineIdx == -1 { // EOF reached for all files.
			break
		}

		output(lines[lesserLineIdx], &prevLine, flags, &printedEmpty)
		updateLines(scnrs, lines, alive, lesserLineIdx)

	}

	if flags.K && flags.R {
		printBlanks(flags.Blanks)
	}

	return nil

}

// printBlanks prints the specified number of blank lines.
func printBlanks(amount *atomic.Int32) {
	blanks := amount.Load()
	for range blanks {
		fmt.Println()
	}
}

// output prints the current line to stdout according to sorting flags.
func output(outputLine string, prevLine *string, flags *flags.Flags, printedEmpty *bool) {

	shouldPrint := true

	if flags.U {
		if outputLine == "" {
			shouldPrint = !*printedEmpty
			*printedEmpty = true
		} else {
			shouldPrint = (outputLine != *prevLine)
			*prevLine = outputLine
		}

	} else if flags.K {
		shouldPrint = (strings.TrimSpace(outputLine) != "")
	}

	if shouldPrint {
		fmt.Println(outputLine)
	}

}

// updateLines updates scanner state after reading one line from a file.
func updateLines(scanners []*bufio.Scanner, lines []string, alive []bool, lesserLineIdx int) {

	if scanners[lesserLineIdx].Scan() {
		lines[lesserLineIdx] = scanners[lesserLineIdx].Text()
		alive[lesserLineIdx] = true
	} else {
		lines[lesserLineIdx] = ""
		alive[lesserLineIdx] = false
	}

}

// getLesserLineIdx compares two lines and returns the index of the lesser one.
func getLesserLineIdx(minIdx int, curIdx int, lines []string, flags *flags.Flags) int {

	if minIdx == -1 {
		return curIdx
	}

	line1 := prepareLine(lines[minIdx], flags)
	line2 := prepareLine(lines[curIdx], flags)

	if flags.K {
		switch {
		case line1 == "" && line2 != "":
			return curIdx
		case line1 != "" && line2 == "":
			return minIdx
		}
	}

	if comparator.CompareLines(line1, line2, flags) > 0 {
		return curIdx
	}
	return minIdx

}

// prepareLine extracts the key field from a line for comparison.
func prepareLine(line string, flags *flags.Flags) string {
	if !flags.K {
		return line
	}
	fields := strings.Fields(line)
	if len(fields) > flags.ClmnToSort {
		return fields[flags.ClmnToSort]

	}
	return ""
}

// cleanup closes all open files and removes temporary chunk files.
func cleanup(files []*os.File, chunks []string) error {

	if err := closeFiles(files); err != nil {
		return fmt.Errorf("failed to close files: %w", err)
	}

	if err := deleteChunks(chunks); err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	return nil
}

// closeFiles safely closes all opened files.
func closeFiles(files []*os.File) error {
	for _, file := range files {
		if err := file.Close(); err != nil {
			return fmt.Errorf("file.Close() returned an error: %w", err)
		}
	}
	return nil
}

// deleteChunks removes temporary chunk files created during sorting.
func deleteChunks(chunks []string) error {
	for _, chunk := range chunks {
		if err := os.Remove(chunk); err != nil {
			return fmt.Errorf("unable to remove file: %w", err)
		}
	}
	return nil
}

// LogFatal prints formatted error messages consistent with GNU sort behavior
func LogFatal(err error, fileName string) {
	if pathErr, ok := err.(*os.PathError); ok {
		switch {
		case errors.Is(pathErr.Err, syscall.EISDIR):
			fmt.Fprintf(os.Stderr, "sort: read failed: %s: Is a directory\n", fileName)
		case errors.Is(pathErr.Err, os.ErrNotExist):
			fmt.Fprintf(os.Stderr, "sort: cannot read: %s: No such file or directory\n", fileName)
		case errors.Is(pathErr.Err, os.ErrPermission):
			fmt.Fprintf(os.Stderr, "sort: cannot read: %s: Permission denied\n", fileName)
		}
	} else {
		fmt.Fprintf(os.Stderr, "sort: fatal error: %v\n", err)
	}
	os.Exit(internalError) // Exit with code 2 for internal errors, matching GNU sort behavior
}

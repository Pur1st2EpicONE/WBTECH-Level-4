// Package master provides functionality for distributing sorting tasks
// to remote worker nodes and collecting results.
package master

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"L4.2/internal/flags"
)

const clientTimeout = 5 * time.Second

var ErrQuorumNotReached = errors.New("quorum not reached") // quorum not reached

// request represents the JSON request sent to worker nodes.
type request struct {
	Lines []string    `json:"lines"`
	Flags flags.Flags `json:"flags"`
}

// response represents the JSON response received from worker nodes.
type response struct {
	Lines []string `json:"lines"`
}

// result stores the outcome of a remote sorting request.
type result struct {
	lines []string
	err   error
}

// RemoteSort distributes the given chunk to multiple nodes, waits for responses,
// and returns the first result that satisfies the quorum.
// Returns ErrQuorumNotReached if the number of successful responses is below quorum.
func RemoteSort(chunk []string, nodes []string, flags *flags.Flags, quorum int) ([]string, error) {

	resChan := make(chan result, len(nodes))

	for _, node := range nodes {
		go func() {
			res, err := sendAndReceive(node, chunk, flags)
			resChan <- result{res, err}
		}()
	}

	success := 0
	var results [][]string

	for range nodes {

		result := <-resChan

		if result.err == nil {
			success++
			results = append(results, result.lines)
		}

		if success >= quorum {
			return results[0], nil
		}

	}

	return nil, ErrQuorumNotReached

}

// sendAndReceive sends the chunk to a single worker node and returns its sorted response.
func sendAndReceive(node string, chunk []string, flags *flags.Flags) ([]string, error) {

	body, err := json.Marshal(request{Lines: chunk, Flags: *flags})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	client := http.Client{Timeout: clientTimeout}

	resp, err := client.Post("http://"+node+"/sort", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send to %s: %w", node, err)
	}
	defer func() { _ = resp.Body.Close() }()

	var res response

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return res.Lines, err

}

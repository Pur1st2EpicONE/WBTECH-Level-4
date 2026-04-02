// Package server provides the HTTP server that acts as a remote sorting worker.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"L4.2/internal/flags"
	"L4.2/internal/sorter"
)

// request represents the incoming JSON payload for the /sort endpoint.
type request struct {
	Lines []string    `json:"lines"`
	Flags flags.Flags `json:"flags"`
}

// response represents the JSON payload returned from the /sort endpoint.
type response struct {
	Lines []string `json:"lines"`
}

// ListenAndServe starts an HTTP server on the given port
// and handles /sort requests for remote sorting.
func ListenAndServe(port string) error {
	http.HandleFunc("/sort", func(w http.ResponseWriter, r *http.Request) { sort(w, r, port) })
	return http.ListenAndServe(":"+port, nil)
}

// sort handles a single /sort request, sorts the provided lines using SortChunk,
// and returns the sorted lines in the response.
func sort(writer http.ResponseWriter, req *http.Request, port string) {

	var request request
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		sorter.LogFatal(fmt.Errorf("failed to decode JSON: %w", err), "")
	}

	fmt.Printf("server %s received chunk, size: %d\n", port, len(request.Lines))

	sorter.SortChunk(request.Lines, &request.Flags)

	if err := json.NewEncoder(writer).Encode(response{Lines: request.Lines}); err != nil {
		sorter.LogFatal(fmt.Errorf("failed to encode JSON: %w", err), "")
	}

}

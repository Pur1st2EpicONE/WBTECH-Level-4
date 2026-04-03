## L4.2

![sort banner](assets/banner.png)

<h3 align="center"> A distributed implementation of the UNIX sort utility in Go, featuring quorum-based result aggregation, concurrency, and large-scale data processing.</h3>

<br>

## Overview

This project extends a classic UNIX-like sort utility into a distributed system.

Each instance of the program can operate in two modes:

- Worker mode — runs an HTTP server and sorts incoming chunks
- Master mode — splits input, distributes work across nodes, and aggregates results

The final result is produced when a quorum (N/2 + 1) of worker nodes successfully process their chunks.

<br>

## Supported Flags

#### Sorting

-k N, --key — Sort by the N-th column (tab-separated).

-n, --numeric-sort — Compare lines according to their numeric value.

-r, --reverse — Reverse the result of comparisons.

-u, --unique — With -c, check for strict ordering; without -c, output only the first of equal lines.

-M, --month-sort — Compare month names (JAN < ... < DEC).

-b, --ignore-leading-blanks — Ignore leading blanks in lines.

-c, --check — Check if input is sorted; do not sort.

-h, --human-numeric-sort — Compare human-readable numbers (e.g., 2K, 1G).


#### Distributed Mode

--nodes — comma-separated list of worker nodes (host:port)

--quorum — minimum number of successful responses required

--serve — run in worker mode

--port — port for worker server (default: 8080)

<br>

## Installation and usage

1) Edit config.yaml to set the number of lines per chunk and number of concurrent workers, if needed.

2) Build the project:

```bash
make
```
3) Start worker nodes (in separate terminals):

```bash
./sort --serve --port 8081
./sort --serve --port 8082
./sort --serve --port 8083
```

4) Run distributed sort:

```bash
./sort --nodes localhost:8081,localhost:8082,localhost:8083 --quorum 2 file_to_sort.txt
```

<br>

## Cool features

* Distributed Sorting with Quorum: chunks are sent to multiple worker nodes, and the result is accepted once the quorum is reached.

* Concurrent Processing: uses goroutines and channels for parallel chunk sorting and asynchronous communication with workers.

* GNU-like Error Handling: mimics GNU sort exit codes and error messages.

<br>

## Testing & Linting

Run tests and ensure code quality:

```bash
make test        # Unit tests
make diff_test   # Differential testing
make lint        # Linting checks
```

The project includes a comprehensive diff-based test suite that compares results with the original UNIX sort. The script automatically builds the binary, starts multiple worker nodes and runs all test cases. 

After running differential tests, a test.log file is generated containing logs from worker nodes. This log demonstrates how chunks are distributed and processed across different servers in parallel.

⚠️ Note: Tests are designed for Linux environments, macOS may produce slightly different results.
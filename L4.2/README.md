## L2.10

![sort banner](assets/banner.png)

<h3 align="center">A simplified implementation of the UNIX sort utility in Go, supporting multiple sorting modes, large file handling, and GNU-style flags.</h3>

<br>

## Supported Flags

-k N, --key — Sort by the N-th column (tab-separated).

-n, --numeric-sort — Compare lines according to their numeric value.

-r, --reverse — Reverse the result of comparisons.

-u, --unique — With -c, check for strict ordering; without -c, output only the first of equal lines.

-M, --month-sort — Compare month names (JAN < ... < DEC).

-b, --ignore-leading-blanks — Ignore leading blanks in lines.

-c, --check — Check if input is sorted; do not sort.

-h, --human-numeric-sort — Compare human-readable numbers (e.g., 2K, 1G).

<br>

## Installation and usage

1) Edit config.yaml to set the number of lines per chunk and number of concurrent workers, if needed.

2) Build the project:

```bash
make
```
3) Run the utility:

```bash
./sort [OPTION]... [FILE]...
```

<br>

## Cool features

* External Sorting: Efficiently handles large files by splitting input into chunks, sorting each chunk, and merging results.

* Concurrent Chunk Processing: Uses Goroutines to sort chunks in parallel.

* GNU-like Error Handling: Mimics GNU sort exit codes and error messages.

* Blazingly Fast: Only about three times slower than real GNU sort on 1,000,000 lines :^)

<br>

## Testing & Linting

Run tests and ensure code quality:

```bash
make test        # Unit tests
make diff_test   # Differential tests comparing output with GNU sort (tested on Linux; results may differ on macOS)
make lint        # Linting checks
```
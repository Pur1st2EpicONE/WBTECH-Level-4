#!/usr/bin/env bash

SORT_OUTPUT="sort_output"
MY_SORT_OUTPUT="my_sort_output"

run_test() {
    
    local file="$1"
    
    sort "$file" > "$SORT_OUTPUT"
    ./sort "$file" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: no flags, one file"
    else
        echo "============================================"
        echo "Test failed: no flags, one file"
        echo "Expected (sort):"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi

    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort "$file" "./assets/test_file_2.txt" > "$SORT_OUTPUT"
    ./sort "$file" "./assets/test_file_2.txt" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: no flags, two files"
    else
        echo "============================================"
        echo "Test failed: no flags, two files"
        echo "Expected (sort):"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi

    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -k 2 "./assets/test_file_2.txt" > "$SORT_OUTPUT"
    ./sort "./assets/test_file_2.txt" -k 2 > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -k 2"
    else
        echo "============================================"
        echo "Test failed: -k 2"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -k 3 "./assets/test_file_2.txt" > "$SORT_OUTPUT"
    ./sort "./assets/test_file_2.txt" -k 3 > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -k 3"
    else
        echo "============================================"
        echo "Test failed: -k 3"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -k 4 "./assets/test_file_2.txt" > "$SORT_OUTPUT"
    ./sort "./assets/test_file_2.txt" -k 4 > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -k 4"
    else
        echo "============================================"
        echo "Test failed: -k 4"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -n "$file" > "$SORT_OUTPUT"
    ./sort -n "$file" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -n"
    else
        echo "============================================"
        echo "Test failed: -n"
        echo "Expected (sort -n):"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -r "$file" > "$SORT_OUTPUT"
    ./sort -r "$file" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -r"
    else
        echo "============================================"
        echo "Test failed: -r"
        echo "Expected (sort):"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"
    
    sort -u "$file" > "$SORT_OUTPUT"
    ./sort -u "$file" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -u"
    else
        echo "============================================"
        echo "Test failed: -u"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -M "$file" > "$SORT_OUTPUT"
    ./sort "$file" -M > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -M"
    else
        echo "============================================"
        echo "Test failed: -M"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -b "$file" > "$SORT_OUTPUT"
    ./sort "$file" -b > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -b"
    else
        echo "============================================"
        echo "Test failed: -b"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -c "$file" > "$SORT_OUTPUT" 2>&1
    ./sort -c "$file" > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -c"
    else
        echo "============================================"
        echo "Test failed: -c"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -h "$file" > "$SORT_OUTPUT" 
    ./sort -h "$file" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -h"
    else
        echo "============================================"
        echo "Test failed: -h"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -rn "$file" > "$SORT_OUTPUT"
    ./sort -rn "$file" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -rn"
    else
        echo "============================================"
        echo "Test failed: -rn"
        echo "Expected (sort):"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -ru "$file" > "$SORT_OUTPUT"
    ./sort -ru "$file" > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -ru"
    else
        echo "============================================"
        echo "Test failed: -ru"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -rM "$file" > "$SORT_OUTPUT"
    ./sort "$file" -rM > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -rM"
    else
        echo "============================================"
        echo "Test failed: -rM"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -rb "$file" > "$SORT_OUTPUT"
    ./sort "$file" -rb > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -rb"
    else
        echo "============================================"
        echo "Test failed: -rb"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -rc "$file" > "$SORT_OUTPUT" 2>&1
    ./sort -rc "$file" > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -rc"
    else
        echo "============================================"
        echo "Test failed: -rc"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -rh "$file" > "$SORT_OUTPUT" 2>&1
    ./sort -rh "$file" > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -rh"
    else
        echo "============================================"
        echo "Test failed: -rh"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -rk 2 "./assets/test_file_2.txt" > "$SORT_OUTPUT"
    ./sort "./assets/test_file_2.txt" -rk 2 > "$MY_SORT_OUTPUT"

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -rk 2"
    else
        echo "============================================"
        echo "Test failed: -rk 2"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -cb "./assets/test_file_2.txt" > "$SORT_OUTPUT" 2>&1
    ./sort "./assets/test_file_2.txt" -cb > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -cb"
    else
        echo "============================================"
        echo "Test failed: -cb"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -cn "./assets/test_file_2.txt" > "$SORT_OUTPUT" 2>&1
    ./sort "./assets/test_file_2.txt" -cn > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -cn"
    else
        echo "============================================"
        echo "Test failed: -cn"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -cu "./assets/test_file_2.txt" > "$SORT_OUTPUT" 2>&1
    ./sort "./assets/test_file_2.txt" -cu > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -cu"
    else
        echo "============================================"
        echo "Test failed: -cu"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -cM "./assets/test_file_2.txt" > "$SORT_OUTPUT" 2>&1
    ./sort "./assets/test_file_2.txt" -cM > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -cM"
    else
        echo "============================================"
        echo "Test failed: -cM"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -ch "./assets/test_file_2.txt" > "$SORT_OUTPUT" 2>&1
    ./sort "./assets/test_file_2.txt" -ch > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -ch"
    else
        echo "============================================"
        echo "Test failed: -ch"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort -ck 2 "./assets/test_file_2.txt" > "$SORT_OUTPUT" 2>&1
    ./sort "./assets/test_file_2.txt" -ck 2 > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: -ck 2"
    else
        echo "============================================"
        echo "Test failed: -ck 2"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort "./internal" > "$SORT_OUTPUT" 2>&1
    ./sort "./internal" > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: Is a directory error"
    else
        echo "============================================"
        echo "Test failed: Is a directory error"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

    sort "NoSuchFile" > "$SORT_OUTPUT" 2>&1
    ./sort "NoSuchFile" > "$MY_SORT_OUTPUT" 2>&1

    if diff -u "$SORT_OUTPUT" "$MY_SORT_OUTPUT"; then
        echo "Test passed: No such file or directory error"
    else
        echo "============================================"
        echo "Test failed: No such file or directory error"
        echo "Expected:"
        cat "$SORT_OUTPUT"
        echo "--------------------------------------------"
        echo "Got:"
        cat "$MY_SORT_OUTPUT"
        echo "============================================"
    fi
    
    rm -f "$SORT_OUTPUT" "$MY_SORT_OUTPUT"

}

go build -o sort ./cmd/sort/main.go

run_test "./assets/test_file_1.txt"

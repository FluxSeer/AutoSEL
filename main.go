package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "the legacy AutoSEL runtime is disabled; use 'go run ./cmd/autosel plan --help'")
	os.Exit(2)
}

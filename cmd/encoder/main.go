package main

import (
	"fmt"
	"os"

	"github.com/VictoriqueMoe/umineko_script_parser/decoder"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: encoder <input.txt> <output.file>\n")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}

	encoded, err := decoder.Encode(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(os.Args[2], encoded, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d bytes -> %d bytes\n", len(data), len(encoded))
}

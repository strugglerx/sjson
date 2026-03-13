package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/strugglerx/sjson"
)

func main() {
	var inPath string
	var outPath string

	flag.StringVar(&inPath, "in", "", "input JSON file, defaults to stdin")
	flag.StringVar(&outPath, "out", "", "output file, defaults to stdout")
	flag.Parse()

	input, err := readInput(inPath)
	if err != nil {
		exitf("read input: %v", err)
	}

	output, err := sjson.StringWithJsonScanToBytesE(json.RawMessage(input))
	if err != nil {
		exitf("expand json string fields: %v", err)
	}

	if err := writeOutput(outPath, output); err != nil {
		exitf("write output: %v", err)
	}
}

func readInput(path string) ([]byte, error) {
	if path == "" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func writeOutput(path string, data []byte) error {
	if path == "" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func exitf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ImmortaL-jsdev/go-learning/go-cli-tools/internal/counter"
)

func main() {
	linesFlag := flag.Bool("l", false, "print line count")
	wordsFlag := flag.Bool("w", false, "print words count")
	bytesFlag := flag.Bool("c", false, "print bytes count")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: wc [-l] [-w] [-c] <filename>")
		os.Exit(1)
	}

	filename := args[0]

	lines, words, bytes, err := counter.CountFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !*linesFlag && !*wordsFlag && !*bytesFlag {
		*linesFlag = true
		*wordsFlag = true
		*bytesFlag = true
	}

	if *linesFlag {
		fmt.Printf("%8d", lines)
	}
	if *wordsFlag {
		fmt.Printf("%8d", words)
	}
	if *bytesFlag {
		fmt.Printf("%8d", bytes)
	}

	fmt.Printf(" %s\n", filename)
}

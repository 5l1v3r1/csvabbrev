package main

import (
	"fmt"
	"io"
	"os"
)

func ProcessArgs() (inflate bool, input io.ReadCloser, output io.WriteCloser) {
	var inputFile *string
	var outputFile *string

	switch len(os.Args) {
	case 1:
		return false, os.Stdin, os.Stdout
	case 2:
		switch os.Args[1] {
		case "-h", "--help":
			printUsage()
			os.Exit(0)
		case "-i", "--inflate":
			inflate = true
		default:
			inputFile = &os.Args[1]
		}
	case 3:
		switch os.Args[1] {
		case "-h", "--help":
			// I use dieUsage() instead of printUsage() because -h shouldn't be used with any other
			// arguments, so their usage of -h was in error.
			dieUsage()
		case "-i", "--inflate":
			inflate = true
			inputFile = &os.Args[2]
		default:
			inputFile = &os.Args[1]
			outputFile = &os.Args[2]
		}
	case 4:
		switch os.Args[1] {
		case "-i", "--inflate":
			inflate = true
			inputFile = &os.Args[2]
			outputFile = &os.Args[3]
		default:
			dieUsage()
		}
	}

	var err error
	if inputFile == nil {
		input = os.Stdin
	} else {
		input, err = os.Open(*inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if outputFile == nil {
		output = os.Stdout
	} else {
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	return
}

func dieUsage() {
	printUsage()
	os.Exit(1)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-h | --help] [-i | --inflate] [file-in [file-out]]\n",
		os.Args[0])
}

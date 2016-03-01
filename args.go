package main

import (
	"fmt"
	"io"
	"os"
)

type Flags struct {
	Inflate bool
	Buffer  bool
	help    bool
}

func ProcessArgs() (flags Flags, input io.ReadCloser, output io.WriteCloser) {
	var inputFile *string
	var outputFile *string
	var args []string

	flags, args = parseFlags()
	if flags.help {
		printUsage()
		if len(os.Args) != 2 {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	switch len(args) {
	case 0:
	case 1:
		inputFile = &args[0]
	case 2:
		inputFile = &args[0]
		outputFile = &os.Args[1]
	default:
		dieUsage()
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
			input.Close()
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	return
}

func parseFlags() (flags Flags, remaining []string) {
	validFlags := map[string]*bool{
		"-i":        &flags.Inflate,
		"--inflate": &flags.Inflate,
		"-h":        &flags.help,
		"--help":    &flags.help,
		"-b":        &flags.Buffer,
		"--buffer":  &flags.Buffer,
	}
	for i, arg := range os.Args[1:] {
		if ptr := validFlags[arg]; ptr == nil {
			return flags, os.Args[i+1:]
		} else {
			if *ptr {
				fmt.Fprintln(os.Stderr, "redundant flag:", arg)
				os.Exit(1)
			}
			*ptr = true
		}
	}
	return flags, []string{}
}

func dieUsage() {
	printUsage()
	os.Exit(1)
}

func printUsage() {
	newlineSpaceCount := len(fmt.Sprintf("Usage: %s ", os.Args[0]))
	newlineSpaces := ""
	for i := 0; i < newlineSpaceCount; i++ {
		newlineSpaces += " "
	}
	fmt.Fprintf(os.Stderr, "Usage: %s [-h | --help] [-i | --inflate] [-b | --buffer]\n"+
		newlineSpaces+"[file-in [file-out]]\n",
		os.Args[0])
}

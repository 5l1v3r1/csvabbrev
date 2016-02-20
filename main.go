package main

import (
	"fmt"
	"os"
)

func main() {
	inflate, input, output := processArgs()
	defer input.Close()
	defer output.Close()

	var err error
	if inflate {
		err = inflateStream(input, output)
	} else {
		err = deflateStream(input, output)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

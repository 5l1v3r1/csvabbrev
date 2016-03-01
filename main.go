package main

import (
	"fmt"
	"os"
	"os/signal"
)

var DieChannel = make(chan struct{}, 1)

func main() {
	flags, input, output := ProcessArgs()

	// On OS X, closing os.Stdin blocks when another Goroutine is reading on stdin.
	if input != os.Stdin {
		defer input.Close()
	}

	defer output.Close()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Fprintln(os.Stderr, "Caught signal. Terminating...")
		close(DieChannel)
	}()

	var err error
	if flags.Inflate {
		err = InflateStream(input, output, flags.Buffer)
	} else {
		err = DeflateStream(input, output, flags.Buffer)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

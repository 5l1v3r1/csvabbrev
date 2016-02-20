package main

import (
	"fmt"
	"os"
	"os/signal"
)

var DieChannel = make(chan struct{}, 1)

func main() {
	inflate, input, output := processArgs()
	defer input.Close()
	defer output.Close()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		close(DieChannel)
	}()

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

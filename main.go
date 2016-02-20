package main

func main() {
	inflate, input, output := processArgs()
	defer input.Close()
	defer output.Close()

	if inflate {
		inflateStream(input, output)
	} else {
		deflateStream(input, output)
	}
}

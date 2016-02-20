package main

import (
	"encoding/csv"
	"io"
)

// ReadCSV reads a CSV stream line by line, passing each parsed record to f.
//
// If a line fails to be read, then this returns a non-nil error.
// When this hits EOF, it will return nil.
//
// If f returns a non-nil value, this will return that same value.
//
// This will return nil if DieChannel is closed.
func ReadCSV(r io.Reader, f func(record []string) error) error {
	stream := readCSVLines(r)
	for {
		select {
		case <-DieChannel:
			return nil
		default:
		}

		select {
		case x, ok := <-stream:
			if !ok {
				return nil
			} else if x.err != nil {
				return x.err
			}
			if err := f(x.record); err != nil {
				return err
			}
		case <-DieChannel:
			return nil
		}
	}
}

// IsCarrotOnly returns true if the given string contains at least one '^' and no other types of
// character.
func IsCarrotOnly(entry string) bool {
	if len(entry) < 1 {
		return false
	}
	for _, ch := range entry {
		if ch != '^' {
			return false
		}
	}
	return true
}

type csvLine struct {
	record []string
	err    error
}

func readCSVLines(r io.Reader) <-chan csvLine {
	res := make(chan csvLine)
	go func() {
		defer close(res)
		r := csv.NewReader(r)
		for {
			record, err := r.Read()
			if record == nil && err == io.EOF {
				return
			}
			res <- csvLine{record, err}
			if err != nil {
				return
			}
		}
	}()
	return res
}

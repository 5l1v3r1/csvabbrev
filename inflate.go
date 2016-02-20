package main

import (
	"encoding/csv"
	"errors"
	"io"
)

func InflateStream(r io.Reader, w io.Writer) error {
	var lastRecord []string
	csvWriter := csv.NewWriter(w)
	return ReadCSV(r, func(record []string) error {
		if err := inflateRecord(record, lastRecord); err != nil {
			return err
		}
		lastRecord = record
		if err := csvWriter.Write(record); err != nil {
			return err
		}
		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			return err
		}
		return nil
	})
}

func inflateRecord(record, lastRecord []string) error {
	for i, x := range record {
		if x == "^" {
			if lastRecord == nil {
				return errors.New("cannot use '^' in line 1")
			} else {
				record[i] = lastRecord[i]
			}
		} else if IsCarrotOnly(x) {
			record[i] = x[1:]
		}
	}
	return nil
}

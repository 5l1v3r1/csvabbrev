package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

func inflateStream(r io.Reader, w io.Writer) error {
	var lastRecord []string
	csvReader := csv.NewReader(r)
	csvWriter := csv.NewWriter(w)
	lineNum := 1
	for {
		record, err := csvReader.Read()
		if err != nil {
			return errors.New("line " + strconv.Itoa(lineNum) + ": " + err.Error())
		}
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
		lineNum++
		select {
		case <-DieChannel:
			return nil
		default:
		}
	}
}

func inflateRecord(record, lastRecord []string) error {
	for i, x := range record {
		if x == "^" {
			if lastRecord == nil {
				return errors.New("cannot use '^' in first entry")
			} else {
				record[i] = lastRecord[i]
			}
		} else if isCarrotOnly(x) {
			record[i] = x[1:]
		}
	}
	return nil
}

func isCarrotOnly(entry string) bool {
	if len(entry) < 2 {
		return false
	}
	for _, ch := range entry {
		if ch != '^' {
			return false
		}
	}
	return true
}
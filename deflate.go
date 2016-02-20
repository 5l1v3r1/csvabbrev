package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

func deflateStream(r io.Reader, w io.Writer) error {
	var lastRecord []string
	csvReader := csv.NewReader(r)
	csvWriter := csv.NewWriter(w)
	lineNum := 1
	for {
		record, err := csvReader.Read()
		if err != nil {
			return errors.New("line " + strconv.Itoa(lineNum) + ": " + err.Error())
		}
		deflated := deflateRecord(record, lastRecord)
		lastRecord = record
		if err := csvWriter.Write(deflated); err != nil {
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

func deflateRecord(record, lastRecord []string) []string {
	if lastRecord == nil {
		return record
	}
	res := make([]string, len(record))
	for i, x := range record {
		if x == lastRecord[i] {
			res[i] = "^"
		} else if isCarrotOnly(x) {
			res[i] = "^" + x
		} else {
			res[i] = x
		}
	}
	return res
}

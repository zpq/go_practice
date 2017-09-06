package main

import (
	"bufio"
	"encoding/csv"
	"os"
)

func LoadFile(filename string) (record [][]string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	r := csv.NewReader(bufio.NewReader(f))

	record, err = r.ReadAll()
	if err != nil {
		return
	}
	return
}

package main

import (
	"bytes"
	"os"
)

func loadJs(filename string) (string, error) {
	result := ""
	f, err := os.Open(filename)
	if err != nil {
		return result, err
	}

	buff := bytes.NewBuffer(nil)
	_, err = buff.ReadFrom(f)
	if err != nil {
		return result, err
	}

	return buff.String(), nil
}

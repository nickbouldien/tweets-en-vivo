package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func Read(reader bufio.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)

	line, err := reader.ReadBytes('\n')

	if err != nil && err != io.EOF {
		// all errors other than the end of file error
		return nil, err
	}
	if err == io.EOF && len(line) == 0 {
		if buf.Len() == 0 {
			return nil, err
		}
	}
	buf.Write(line)

	return buf.Bytes(), nil
}

// PrettyPrint is a helper function to print the data to the terminal with some formatting
func PrettyPrint(data interface{}) {
	s, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		_ = fmt.Errorf("error creating the data to print: %v", err)
	}

	fmt.Println(string(s))
}

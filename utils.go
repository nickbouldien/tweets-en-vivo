package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
func PrettyPrint(data []byte) {
	// TODO - clean this up
	var rules bytes.Buffer
	if err := json.Indent(&rules, data,"","\t"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(rules.Bytes()))
}


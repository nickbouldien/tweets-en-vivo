package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// CloseFile is a helper function to close a file. It will os.Exit(1) if there is an error closing the file
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// PrettyPrint is a helper function to print the data to the terminal with some formatting
func PrettyPrint(data interface{}) {
	s, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		_ = fmt.Errorf("error creating the data to print: %v", err)
	}

	fmt.Println(string(s))
}

// PrettyPrintByteSlice is a helper function to print the slice of bytes to the terminal with some formatting
func PrettyPrintByteSlice(data []byte) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "\t"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(buf.Bytes()))
}

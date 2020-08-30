package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func (r *StreamResponseBodyReader) Read() ([]byte, error) {
	r.buf.Truncate(0)

	for {
		line, err := r.reader.ReadBytes('\n')

		if len(line) == 0 {
			fmt.Println("len(line) == 0")
			continue
		}

		if err != nil && err != io.EOF {
			// all errors other than the end of file error
			_ = fmt.Errorf("read error: %v", err)
			return nil, err
		}

		if err == io.EOF && len(line) == 0 {
			_ = fmt.Errorf("io.EOF && len(line): %v", err)

			if r.buf.Len() == 0 {
				_ = fmt.Errorf("buf.Len() : %v", err)
				return nil, err
			}
			fmt.Println("breaking")
			break
		}

		if bytes.HasSuffix(line, []byte("\r\n")) {
			r.buf.Write(bytes.TrimRight(line, "\r\n"))
			break
		}

		r.buf.Write(line)
	}

	return r.buf.Bytes(), nil
}

//func Read(reader bufio.Reader) ([]byte, error) {
//	buf := new(bytes.Buffer)
//
//	line, err := reader.ReadBytes('\n')
//	fmt.Println("called read: ", line)
//
//	if err != nil && err != io.EOF {
//		// all errors other than the end of file error
//		_ = fmt.Errorf("read error: %v", err)
//		return nil, err
//	}
//
//	if err == io.EOF && len(line) == 0 {
//		_ = fmt.Errorf("io.EOF && len(line): %v", err)
//
//		if buf.Len() == 0 {
//			_ = fmt.Errorf("buf.Len() : %v", err)
//			return nil, err
//		}
//		fmt.Println("continuing")
//	}
//
//	if bytes.HasSuffix(line, []byte("\r\n")) {
//		buf.Write(bytes.TrimRight(line, "\r\n"))
//	} else {
//		buf.Write(line)
//	}
//
//	return buf.Bytes(), nil
//}

// PrettyPrint is a helper function to print the data to the terminal with some formatting
func PrettyPrint(data interface{}) {
	s, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		_ = fmt.Errorf("error creating the data to print: %v", err)
	}

	fmt.Println(string(s))
}

func PrettyPrintByteSlice(data []byte) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "\t"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(buf.Bytes()))
}

func CloseFile(file *os.File) {
	fmt.Println("closing file: ", file.Name())
	err := file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

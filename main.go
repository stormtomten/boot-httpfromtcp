package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Unable to open %s:  %v", inputFilePath, err)
	}
	defer file.Close()
	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)
		line := ""

		for {
			buf := make([]byte, 8, 8)
			n, err := f.Read(buf)
			if err != nil {
				if line != "" {
					ch <- line
					line = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				// fmt.Printf("error: %s\n", err)
				break
			}
			str := string(buf[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- line + parts[i]
				line = ""
			}
			line += parts[len(parts)-1]
		}
	}()

	return ch
}

package main

import (
	"boot-httpfromtcp/internal/request"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error litening for TCP traffic on %s: %s", port, err.Error())
	}
	defer l.Close()

	fmt.Printf("Listening for TCP traffic on: %s\n", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Failed to establish connection from %s: %s", conn.RemoteAddr(), err.Error())
		}

		fmt.Printf("Connection established from: %s\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error reading: %s", err.Error())
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		line := ""

		for {
			buf := make([]byte, 8, 8)
			n, err := f.Read(buf)
			if err != nil {
				if line != "" {
					lines <- line
					line = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err)
				return
			}
			str := string(buf[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- line + parts[i]
				line = ""
			}
			line += parts[len(parts)-1]
		}
	}()

	return lines
}

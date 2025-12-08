package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type requestStatus int

const (
	initialized requestStatus = iota
	done
)

type Request struct {
	RequestLine RequestLine
	state       requestStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	crlf       = "\r\n"
	bufferSize = 8
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := Request{
		state: initialized,
	}

	for req.state != done {

		if len(buf) <= readToIndex {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if errors.Is(err, io.EOF) {
			req.state = done
		} else if err != nil {
			return nil, err
		}

		if 0 < n {
			readToIndex += n
			nParsed, err := req.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}
			copy(buf, buf[nParsed:])
			readToIndex -= nParsed
		}

	}
	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case initialized:
		line, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *line
		r.state = done
		return n, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: Unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineString := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineString)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid Request Line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("Incorrect Method Format: %s", method)
		}
	}

	target := parts[1]
	if len(target) == 0 || target[0] != '/' {
		return nil, fmt.Errorf("Invalid request target")
	}

	ver, ok := strings.CutPrefix(parts[2], "HTTP/")
	if !ok {
		return nil, fmt.Errorf("Invalid Version Format")
	}
	if ver != "1.1" {
		return nil, fmt.Errorf("Unsupported HTTP version")
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   ver,
	}, nil
}

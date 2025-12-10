package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

var reasonPhrase = map[StatusCode]string{
	StatusOK:            "OK",
	StatusBadRequest:    "Bad Request",
	StatusInternalError: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := reasonPhrase[statusCode]
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason)

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	return nil
}

package response

import (
	"fmt"
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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writeStatus != writeStatus {
		return fmt.Errorf("error: Wrong write order: %d", w.writeStatus)
	}
	reason := reasonPhrase[statusCode]
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason)

	_, err := w.writer.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	w.writeStatus = writeHeaders
	return nil
}

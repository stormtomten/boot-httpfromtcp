package response

import (
	"boot-httpfromtcp/internal/headers"
	"fmt"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writeStatus != writeHeaders {
		return fmt.Errorf("error: Wrong write order: %d", w.writeStatus)
	}
	for key, val := range h {
		header := fmt.Sprintf("%s: %s\r\n", key, val)
		_, err := w.writer.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	w.writeStatus = writeBody
	return nil
}

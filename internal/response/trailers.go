package response

import (
	"boot-httpfromtcp/internal/headers"
	"fmt"
)

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writeStatus != writeTrailers {
		return fmt.Errorf("can't write trailers in writeStatus: %d", w.writeStatus)
	}

	defer func() { w.writeStatus = writeDone }()

	for key, val := range h {
		header := fmt.Sprintf("%s: %s\r\n", key, val)
		fmt.Printf("Trailers: %s\n", header)
		_, err := w.writer.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))

	return err
}

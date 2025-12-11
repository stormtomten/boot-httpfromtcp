package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writeStatus != writeBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writeStatus)
	}
	chuckSize := len(p)
	total := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chuckSize)
	if err != nil {
		return total + n, err
	}
	total += n
	n, err = w.writer.Write(p) // data
	if err != nil {
		return total + n, err
	}
	total += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return total + n, err
	}
	return total, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writeStatus != writeBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writeStatus)
	}

	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	w.writeStatus = writeTrailers
	return n, nil
}

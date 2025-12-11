package response

import (
	"fmt"
	"io"
	"net"
)

type writerState int

const (
	writeStatus writerState = iota
	writeHeaders
	writeBody
	writeTrailers
	writeDone
)

type Writer struct {
	writer      io.Writer
	writeStatus writerState
}

func NewWriter(conn net.Conn) *Writer {
	w := io.Writer(conn)

	return &Writer{
		writer:      w,
		writeStatus: writeStatus,
	}
}

func (w *Writer) WriteBody(p []byte) error {
	if w.writeStatus != writeBody {
		return fmt.Errorf("cannot write body in state %d", w.writeStatus)
	}
	if _, err := w.writer.Write(p); err != nil {
		return err
	}
	w.writeStatus = writeDone
	return nil
}

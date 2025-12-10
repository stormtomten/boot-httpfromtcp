package response

import (
	"io"
	"net"
)

type writerState int

const (
	writeStatus writerState = iota
	writeHeaders
	writeBody
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
	if _, err := w.writer.Write(p); err != nil {
		return err
	}
	w.writeStatus = writeDone
	return nil
}

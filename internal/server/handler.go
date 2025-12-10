package server

import (
	"boot-httpfromtcp/internal/request"
	"boot-httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	msgbytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(msgbytes))
	response.WriteHeaders(w, headers)
	w.Write(msgbytes)
}

package server

import (
	"boot-httpfromtcp/internal/request"
	"boot-httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

/*
func (he *HandlerError) Write(w *response.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	msgbytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(msgbytes))
	response.WriteHeaders(w, headers)
	w.write.Write(msgbytes)
}
*/

package main

import (
	"boot-httpfromtcp/internal/request"
	"boot-httpfromtcp/internal/response"
	"strings"
)

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	message := []byte(strings.Join([]string{
		"<html>",
		"  <head>",
		"    <title>200 OK</title>",
		"  </head>",
		"  <body>",
		"    <h1>Success!</h1>",
		"    <p>Your request was an absolute banger.</p>",
		"  </body>",
		"</html>",
	}, "\n"))
	responesHeaders := response.GetDefaultHeaders(len(message))
	responesHeaders.Override("Content-Type", "text/html")
	w.WriteHeaders(responesHeaders)
	w.WriteBody(message)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	message := []byte(strings.Join([]string{
		"<html>",
		"  <head>",
		"    <title>400 Bad Request</title>",
		"  </head>",
		"  <body>",
		"    <h1>Bad Request</h1>",
		"    <p>Your request honestly kinda sucked.</p>",
		"  </body>",
		"</html>",
	}, "\n"))

	responesHeaders := response.GetDefaultHeaders(len(message))
	responesHeaders.Override("Content-Type", "text/html")
	w.WriteHeaders(responesHeaders)
	w.WriteBody(message)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalError)
	message := []byte(strings.Join([]string{
		"<html>",
		"  <head>",
		"    <title>500 Internal Server Error</title>",
		"  </head>",
		"  <body>",
		"    <h1>Internal Server Error Request</h1>",
		"    <p>Okay, you know what? This one is on me.</p>",
		"  </body>",
		"</html>",
	}, "\n"))
	responesHeaders := response.GetDefaultHeaders(len(message))
	responesHeaders.Override("Content-Type", "text/html")
	w.WriteHeaders(responesHeaders)
	w.WriteBody(message)
}

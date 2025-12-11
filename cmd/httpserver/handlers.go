package main

import (
	"boot-httpfromtcp/internal/headers"
	"boot-httpfromtcp/internal/request"
	"boot-httpfromtcp/internal/response"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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

const (
	TrailerHash   = "X-Content-SHA256"
	TrailerLength = "X-Content-Length"
)

func handlerChunk(w *response.Writer, req *request.Request) {
	cutTarget := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	path := fmt.Sprintf("https://httpbin.org/%s", cutTarget)
	buf := make([]byte, 1024)

	resp, err := http.Get(path)
	if err != nil {
		w.WriteStatusLine(response.StatusInternalError)
		return
	}

	w.WriteStatusLine(response.StatusCode(resp.StatusCode))
	h := headers.NewHeaders()
	for key, val := range resp.Header {
		if strings.ToLower(key) != "content-length" {
			h.Set(key, strings.Join(val, ", "))
		}
	}

	h.Override("Transfer-Encoding", "chunked")
	h.Set("Trailer", TrailerHash+", "+TrailerLength)
	w.WriteHeaders(h)

	store := make([]byte, 0, 1024)
	for {
		n, err := resp.Body.Read(buf)
		fmt.Printf("Got %d from target\n", n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, err := w.WriteChunkedBodyDone()
				if err != nil {
					fmt.Printf("write close error")
				}
				fmt.Printf("Wrote close to client\n")
			}
			break
		}
		if 0 < n {
			x, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				fmt.Printf("write error")
			}
			store = append(store, buf[:n]...)
			fmt.Printf("Wrote %d to client\n", x)
		}
	}
	sha := sha256.Sum256(store)
	trailer := headers.NewHeaders()
	trailer.Set(TrailerHash, fmt.Sprintf("%x", sha))
	trailer.Set(TrailerLength, fmt.Sprintf("%d", len(store)))
	w.WriteTrailers(trailer)

	return
}

const assetPath = "./assets/vim.mp4"

func handlerVideo(w *response.Writer, _ *request.Request) {
	asset, err := os.ReadFile(assetPath)
	if err != nil {
		fmt.Printf("Failed to read file %s: %s", assetPath, err.Error())
		w.WriteStatusLine(response.StatusInternalError)
		return
	}

	w.WriteStatusLine(response.StatusOK)
	head := headers.NewHeaders()
	head.Set("Content-Type", "video/mp4")
	w.WriteHeaders(head)
	w.WriteBody(asset)
}

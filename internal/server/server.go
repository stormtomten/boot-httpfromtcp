package server

import (
	"boot-httpfromtcp/internal/request"
	"boot-httpfromtcp/internal/response"
	//"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Addr     net.Addr
	Listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, h Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to bind to port: %d", port)
	}

	s := &Server{
		Addr:     l.Addr(),
		Listener: l,
		handler:  h,
	}
	s.closed.Store(false)

	go s.listen()

	return s, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("accept error: %s", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("conn error: failed to  parse request: %s\n", err.Error())
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	s.handler(w, req)
	return
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.Listener.Close() != nil {
		return s.Listener.Close()
	}
	return nil
}

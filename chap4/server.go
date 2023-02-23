package chap4

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Monitor struct {
	logger *log.Logger
}

func (m Monitor) Write(p []byte) (int, error) {
	return len(p), m.logger.Output(2, string(p))
}

type Server struct {
	addr        string
	listener    net.Listener
	monitor     Monitor
	payloadSize uint
}

type serverOpts func(*Server) error

func RunServer() error {
	s, err := NewServer()
	if err != nil {
		return err
	}
	return s.Run()
}

func NewServer(opts ...serverOpts) (*Server, error) {
	s := &Server{
		addr:        "127.0.0.1:3000",
		monitor:     Monitor{logger: log.New(os.Stdout, "monitor: ", log.Llongfile)},
		payloadSize: 1 << 24, // 16 MB
	}
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return s, fmt.Errorf("failed to apply option with error: %w", err)
		}
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return s, fmt.Errorf("failed to start listener on addr %q with error: %w", s.addr, err)
	}
	s.listener = listener

	s.monitor.logger.Printf("listening on %q", s.listener.Addr().String())
	return s, nil
}

func WithLogger(l Monitor) serverOpts {
	return func(s *Server) error {
		s.monitor = l
		return nil
	}
}

func WithAddr(a string) serverOpts {
	return func(s *Server) error {
		s.addr = a
		return nil
	}
}

func WithPayloadSize(sz uint) serverOpts {
	return func(s *Server) error {
		s.payloadSize = sz
		return nil
	}
}

func (s *Server) Run() error {
	payloads := []string{"hello", "there", "general", "kenobi"}
	conn, err := s.listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept listener: %w", err)
	}
	defer conn.Close()
	w := io.MultiWriter(conn, s.monitor)

	for _, p := range payloads {
		// make a random payload
		payload := []byte(p)
		_, err = w.Write(payload)
		if err != nil {
			return fmt.Errorf("failed to write payload: %w", err)
		}
	}
	return nil
}

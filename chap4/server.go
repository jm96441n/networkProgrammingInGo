package chap4

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
)

type Server struct {
	addr        string
	listener    net.Listener
	logger      *log.Logger
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
		logger:      log.New(os.Stdout, "", log.Llongfile),
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

	s.logger.Printf("listening on %q", s.listener.Addr().String())
	return s, nil
}

func WithLogger(l *log.Logger) serverOpts {
	return func(s *Server) error {
		s.logger = l
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
	// make a random payload
	payload := make([]byte, s.payloadSize)
	_, err := rand.Read(payload)
	if err != nil {
		return fmt.Errorf("failed to generate random payload with err: %w", err)
	}
	s.logger.Print("created payload")
	conn, err := s.listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept listener: %w", err)
	}
	defer conn.Close()

	s.logger.Print("writing payload")
	s.logger.Printf("payload: %v", payload)
	_, err = conn.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}
	s.logger.Print("wrote payload")
	return nil
}

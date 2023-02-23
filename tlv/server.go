package tlv

import (
	"log"
	"net"
	"os"
)

type Server struct {
	listener net.Listener
	payloads []Payload
}

func RunServer() error {
	return NewServer().Run()
}

func NewServer() Server {
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Printf("error starting server: %s", err)
		os.Exit(1)
	}
	b1 := Binary("Clear is better than clever")
	b2 := Binary("Don't panic")
	s1 := String("errors are values")
	return Server{
		listener: listener,
		payloads: []Payload{&b1, &b2, &s1},
	}
}

func (s Server) Run() error {
	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}

	defer conn.Close()

	for _, p := range s.payloads {
		_, err := p.WriteTo(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

package chap3

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

type Server struct {
	addr     string
	output   io.Writer
	listener net.Listener
}
type serverOption func(*Server)

func RunServer() {
	NewServer().RunServer()
}

func withServerAddr(addr string) serverOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func withServerOutput(output io.Writer) serverOption {
	return func(s *Server) {
		s.output = output
	}
}

func NewServer(opts ...serverOption) *Server {
	s := &Server{
		addr:   "127.0.0.1:3000",
		output: os.Stdout,
	}
	for _, opt := range opts {
		opt(s)
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logAndDie(err)
	}
	s.listener = listener
	fmt.Printf("listening on %s", listener.Addr().String())
	return s
}

func (c *Server) RunServer() {
	defer c.listener.Close()
	wg := &sync.WaitGroup{}
	conn, err := c.listener.Accept()
	if err != nil {
		c.logAndDie(err)
	}
	wg.Add(1)
	go c.HandleConn(conn, wg)
	wg.Wait()
}

func (c *Server) HandleConn(conn io.WriteCloser, wg *sync.WaitGroup) {
	defer conn.Close()
	for i := 0; i <= 10; i++ {
		msg := fmt.Sprintf("msg num %d", i)
		conn.Write([]byte(msg))
		fmt.Fprintf(c.output, "Sending: %s\n", msg)
	}
	wg.Done()
}

func (s *Server) logAndDie(err error) {
	fmt.Fprint(s.output, err)
	os.Exit(1)
}

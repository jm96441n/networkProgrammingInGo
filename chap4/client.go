package chap4

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	logger *log.Logger
	addr   string
}

type clientOpt func(*Client)

func RunClient() error {
	return NewClient().Run()
}

func NewClient(opts ...clientOpt) *Client {
	c := &Client{
		addr:   "127.0.0.1:3000",
		logger: log.New(os.Stdout, "", log.Llongfile),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithClientLogger(l *log.Logger) clientOpt {
	return func(c *Client) {
		c.logger = l
	}
}

func WithClientAddr(a string) clientOpt {
	return func(c *Client) {
		c.addr = a
	}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("failed to dial server at %q: %w", c.addr, err)
	}
	buf := make([]byte, 25)
	results := make([]byte, 0, 0)
	for {
		n, err := conn.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read from conn: %w", err)
		}
		c.logger.Print(buf)
		results = append(results, buf[:n]...)

	}
	c.logger.Print(results)
	return nil
}

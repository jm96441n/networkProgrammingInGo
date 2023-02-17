package chap3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	output io.Writer
	addr   string
}

type clientOption func(*Client)

func RunClient() {
	NewClient().RunClient()
}

func withAddr(addr string) clientOption {
	return func(c *Client) {
		c.addr = addr
	}
}

func withOutput(output io.Writer) clientOption {
	return func(c *Client) {
		c.output = output
	}
}

func NewClient(opts ...clientOption) *Client {
	c := &Client{
		addr:   "127.0.0.1:",
		output: os.Stdout,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) RunClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", "127.0.0.1:3000")
	if err != nil {
		c.logAndDie(err)
	}
	defer conn.Close()
	done := make(chan struct{})

	for {
		select {
		case <-done:
			fmt.Fprintln(c.output, "donezo")
			return
		default:
			go c.HandleConn(conn, done)
		}
	}
}

func (s *Client) HandleConn(conn net.Conn, done chan struct{}) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if errors.Is(err, io.EOF) {
			done <- struct{}{}
			return
		}
		if err != nil {
			s.logAndDie(err)
		}

		fmt.Fprintf(s.output, "received: %q\n", buf[:n])
	}
}

func (c *Client) logAndDie(err error) {
	fmt.Fprint(c.output, err)
	os.Exit(1)
}

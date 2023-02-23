package chap4

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	monitor Monitor
	addr    string
}

type clientOpt func(*Client)

func RunClient() error {
	return NewClient().Run()
}

func NewClient(opts ...clientOpt) *Client {
	c := &Client{
		addr:    "127.0.0.1:3000",
		monitor: Monitor{logger: log.New(os.Stdout, "monitor: ", log.Llongfile)},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithClientLogger(l Monitor) clientOpt {
	return func(c *Client) {
		c.monitor = l
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
	defer conn.Close()

	r := io.TeeReader(conn, c.monitor)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var words []string
	for scanner.Scan() {
		msg := scanner.Text()
		words = append(words, msg)
	}
	err = scanner.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RunNoScanner() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("failed to dial server at %q: %w", c.addr, err)
	}
	defer conn.Close()
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
		c.monitor.logger.Print(buf)
		results = append(results, buf[:n]...)

	}
	c.monitor.logger.Print(results)
	return nil
}

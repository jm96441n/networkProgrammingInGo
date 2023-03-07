package unixstreaming

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/log"
)

func StreamingEchoServer(ctx context.Context, network, addr string) (net.Addr, error) {
	s, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	err = os.Chmod(addr, os.ModeSocket|0666)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := s.Accept()
			if err != nil {
				return
			}

			go func() {
				defer func() { conn.Close() }()
				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					log.Info(fmt.Sprintf("server received: %q", string(buf[:n])))
					_, err = conn.Write(buf[:n])
					if err != nil {
						return
					}
				}
			}()
		}
	}()

	log.Info(fmt.Sprintf("listening on %q ....", s.Addr()))
	<-ctx.Done()
	s.Close()

	return s.Addr(), nil
}

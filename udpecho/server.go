package udpecho

import (
	"context"
	"fmt"
	"net"

	"github.com/charmbracelet/log"
)

func RunServer(ctx context.Context, addr string, numMsgs int, cancel context.CancelFunc) error {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("binding to udp %s: %w", addr, err)
	}

	go func(cancelFn context.CancelFunc) {
		buf := make([]byte, 1024)
		i := 0
		for {
			n, clientAddr, err := s.ReadFrom(buf) // client to server
			if err != nil {
				log.Error(err)
				return
			}

			log.Info(fmt.Sprintf("Received %q from %q", buf[:n], clientAddr.String()))
			_, err = s.WriteTo(buf[:n], clientAddr)
			if err != nil {
				log.Error(err)
				return
			}
			i += 1
			if numMsgs == i {
				log.Info("Finished Receiving")
				cancelFn()
			}
		}
	}(cancel)

	<-ctx.Done()
	_ = s.Close()

	return nil
}

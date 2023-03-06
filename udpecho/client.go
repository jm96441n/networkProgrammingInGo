package udpecho

import (
	"context"
	"fmt"
	"net"

	"github.com/charmbracelet/log"
)

func RunClient(ctx context.Context, serverAddr, clientAddr string, numMsgs int, cancel context.CancelFunc) error {
	client, err := net.ListenPacket("udp", clientAddr)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()
	server, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return err
	}
	for i := 0; i < numMsgs; i++ {
		msg := []byte(fmt.Sprintf("ping %d", i))
		_, err = client.WriteTo(msg, server)
		if err != nil {
			return err
		}

		buf := make([]byte, 1024)
		n, _, err := client.ReadFrom(buf)
		if err != nil {
			return err
		}

		log.Info(fmt.Sprintf("Received %q from %q", buf[:n], server.String()))
	}
	return nil
}

package echo

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func StreamingClient(addr string) error {
	conn, err := net.Dial("unix", addr)
	if err != nil {
		return err
	}

	defer func() { conn.Close() }()

	msg := []byte("ping")

	for i := 0; i < 3; i++ {
		log.Info("client sending ping")
		_, err := conn.Write(msg)
		if err != nil {
			return err
		}
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("client received %q", string(buf[:n])))
		time.Sleep(1 * time.Second)
	}
	return nil
}

func DatagramClient(addr string, network string) error {
	client, err := net.ListenPacket("unixgram", network)
	if err != nil {
		return err
	}

	defer func() { client.Close() }()

	err = os.Chmod(network, os.ModeSocket|0666)
	if err != nil {
		return err
	}

	msg := []byte("ping")

	sAddr, err := net.ResolveUnixAddr("unixgram", addr)
	if err != nil {
		log.Error(fmt.Sprintf("unable to resolve: %v", err))
		return err
	}

	for i := 0; i < 3; i++ {
		log.Info("client sending ping")
		_, err := client.WriteTo(msg, sAddr)
		if err != nil {
			return err
		}
		buf := make([]byte, 1024)

		n, _, err := client.ReadFrom(buf)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("client received %q", string(buf[:n])))
		time.Sleep(1 * time.Second)
	}
	return nil
}

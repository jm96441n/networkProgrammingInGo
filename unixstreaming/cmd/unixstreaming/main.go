package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jm96441n/networkProgrammingInGo/unixstreaming"
)

func main() {
	var (
		clientMode bool
		serverMode bool
		clientAddr string
	)
	flag.BoolVar(&clientMode, "c", false, "run in client mode")
	flag.StringVar(&clientAddr, "a", "", "addr for client to reach server")
	flag.BoolVar(&serverMode, "s", true, "run in client mode")
	flag.Parse()

	if clientMode {

		err := runClient(clientAddr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := runServer()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runServer() error {
	dir, err := os.MkdirTemp("", "echo_unix")
	if err != nil {
		return err
	}

	defer func() {
		rErr := os.RemoveAll(dir)
		if rErr != nil {
			log.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		cancel()
	}()

	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))

	_, err = unixstreaming.StreamingEchoServer(ctx, "unix", socket)
	if err != nil {
		return err
	}

	return nil
}

func runClient(addr string) error {
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

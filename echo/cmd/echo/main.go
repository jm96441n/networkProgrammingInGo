package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/jm96441n/networkProgrammingInGo/echo"
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
	// return runStreamingServer()
	return runDatagramServer()
}

func runStreamingServer() error {
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

	_, err = echo.StreamingEchoServer(ctx, "unix", socket)
	if err != nil {
		return err
	}

	return nil
}

func runDatagramServer() error {
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

	err = echo.DatagramEchoServer(ctx, "unixgram", socket)
	if err != nil {
		return err
	}

	return nil
}

func runClient(addr string) error {
	// return echo.StreamingClient(addr)
	return runDatagramClient(addr)
}

func runDatagramClient(addr string) error {
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
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	return echo.DatagramClient(addr, socket)
}

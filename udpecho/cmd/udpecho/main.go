package main

import (
	"context"
	"flag"

	"github.com/charmbracelet/log"
	"github.com/jm96441n/networkProgrammingInGo/udpecho"
)

func main() {
	var isClient bool
	var isServer bool
	flag.BoolVar(&isClient, "c", false, "Set to run in client mode")
	flag.BoolVar(&isServer, "s", false, "Set to run in server mode")
	flag.Parse()

	if isClient && isServer {
		log.Fatal("Set either client OR server mode, not both")
	}

	ctx, cancel := context.WithCancel(context.Background())
	serverAddr := "127.0.0.1:3000"
	clientAddr := "127.0.0.1:3001"

	switch {
	case isClient:
		err := udpecho.RunClient(ctx, serverAddr, clientAddr, 4, cancel)
		if err != nil {
			log.Fatal(err)
		}
	case isServer:
		err := udpecho.RunServer(ctx, serverAddr, 4, cancel)
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("You must pass either client OR server mod flag")
	}
}

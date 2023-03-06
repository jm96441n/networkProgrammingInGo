package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/jm96441n/networkProgrammingInGo/tftp"
)

func main() {
	addr := "127.0.0.1:3000"
	payload, err := os.ReadFile("./kitten-large.png")
	if err != nil {
		log.Fatal(err)
	}

	server, err := tftp.NewServer(payload)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe(addr))
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jm96441n/networkProgrammingInGo/chap4"
)

func main() {
	var clientMode bool
	var serverMode bool
	flag.BoolVar(&clientMode, "c", false, "run in client mode")
	flag.BoolVar(&serverMode, "s", false, "run in server mode")

	flag.Parse()
	if clientMode {
		err := chap4.RunClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running client: %s", err)
		}
		os.Exit(0)
	}
	err := chap4.RunServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running server: %s", err)
	}
	os.Exit(0)
}

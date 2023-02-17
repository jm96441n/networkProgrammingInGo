package main

import (
	"flag"
	"fmt"

	"github.com/jm96441n/networkProgrammingInGo/chap3"
)

func main() {
	var server, client bool

	flag.BoolVar(&server, "s", false, "use this to run in server mode")
	flag.BoolVar(&client, "c", false, "use this to run in client mode")
	flag.Parse()
	fmt.Println(client)
	fmt.Println(server)

	if server {
		fmt.Println("running server")
		chap3.RunServer()
	} else {
		chap3.RunClient()
	}
}

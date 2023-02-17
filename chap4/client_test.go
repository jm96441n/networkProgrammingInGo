package chap4_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/jm96441n/networkProgrammingInGo/chap4"
)

func TestClientReadsAllMessage(t *testing.T) {
	addr := "127.0.0.1:3002"
	logBuffer := bytes.NewBufferString("")
	logger := log.New(logBuffer, "", log.Lmsgprefix)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	client := chap4.NewClient(chap4.WithClientLogger(logger), chap4.WithClientAddr(addr))
	errs := make(chan error, 1)
	go func() {
		err := client.Run()
		errs <- err
	}()
	conn, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}

	payload := make([]byte, 55)
	_, err = rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Write(payload)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()

	err = <-errs
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(logBuffer.String(), fmt.Sprintf("%v", payload)) {
		t.Errorf("expected %s to contain %v", logBuffer.String(), payload)
	}
}

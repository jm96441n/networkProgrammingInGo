package chap4_test

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"testing"

	"github.com/jm96441n/networkProgrammingInGo/chap4"
)

func TestServerReturnsWritesBackAPayloadGivenAParticularSize(t *testing.T) {
	//	logBuf := bytes.NewBuffer([]byte{})
	logger := log.New(os.Stdout, "", log.Lshortfile)
	addr := "127.0.0.1:3001"
	var payloadSize uint = 50
	s, err := chap4.NewServer(chap4.WithAddr(addr), chap4.WithLogger(logger), chap4.WithPayloadSize(payloadSize))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("starting server")
	go func() {
		err := s.Run()
		if err != nil {
			panic(err)
		}
	}()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("starting to read")
	resp := make([]byte, payloadSize, payloadSize)
	readSize := 0
	for {
		n, err := conn.Read(resp)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		readSize += n
		t.Logf("read %d bytes", n)

	}
	conn.Close()
	if readSize != int(payloadSize) {
		t.Errorf("expected %d to be read, got %d", payloadSize, readSize)
	}
}

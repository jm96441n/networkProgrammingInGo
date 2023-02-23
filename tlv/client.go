package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

type Client struct {
	addr         string
	payloadCount int
}

func RunClient() error {
	return NewClient().Run()
}

func NewClient() Client {
	return Client{
		addr:         "127.0.0.1:3000",
		payloadCount: 3,
	}
}

func (c Client) Run() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	defer conn.Close()

	for i := 0; i < c.payloadCount; i++ {
		actual, err := decode(conn)
		if err != nil {
			return err
		}
		log.Printf("[%T] %[1]q", actual)
	}
	return nil
}

func decode(r io.Reader) (Payload, error) {
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload

	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, errors.New("unknown type")

	}

	_, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, err
	}
	return payload, nil
}

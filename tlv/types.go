package tlv

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayloadSize uint32 = 10 << 20 // 10MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (m *Binary) Bytes() []byte  { return *m }
func (m *Binary) String() string { return string(*m) }

func (m *Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType) // 1-byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(*m))) // 4-byte size

	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(*m) // payload

	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ) // 1-byte type
	if err != nil {
		return 0, err
	}

	if typ != BinaryType {
		return 0, errors.New("invalid Binary")
	}

	var n int64 = 1

	var sz uint32
	err = binary.Read(r, binary.BigEndian, &sz) // 4-byte size
	if err != nil {
		return n, err
	}

	n += 4

	if sz > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}
	*m = make([]byte, sz)
	o, err := r.Read(*m) // payload

	return n + int64(o), err
}

type String string

func (m String) Bytes() []byte  { return []byte(m) }
func (m String) String() string { return string(m) }

func (m String) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType) // 1-byte type
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4-byte size
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(m.Bytes())

	return n + int64(o), err
}

func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8

	err := binary.Read(r, binary.BigEndian, &typ) // 1-byte size
	if err != nil {
		return 0, err
	}

	if typ != StringType {
		return 0, errors.New("invalid String")
	}

	var n int64 = 1

	var sz uint32
	err = binary.Read(r, binary.BigEndian, &sz) // 4-byte size
	if err != nil {
		return n, err
	}

	n += 4

	if sz > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	buf := make([]byte, sz)
	o, err := r.Read(buf)
	if err != nil {
		return n, err
	}

	*m = String(buf)

	return n + int64(o), nil
}

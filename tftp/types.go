package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

const (
	DatagramSize = 516 // max supported datagram size
	HeaderSize   = 4
	BlockSize    = DatagramSize - HeaderSize // the DatagramSize minus the 4-byte header
)

type OpCode uint16

const (
	OpRRQ OpCode = iota + 1 // ReadReQuest
	_                       // no WriteReQuest support
	OpData
	OpAck
	OpErr
)

type ErrCode uint16

const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrFileExists
	ErrNoUser
)

type ReadReq struct {
	Filename string
	Mode     string
}

// not used by the server, but a client would make use of this
func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if q.Mode != "" {
		mode = q.Mode
	}

	// operation code + filename + 0 byte + mode + 0 byte
	cap := 2 + 2 + len(q.Filename) + 1 + len(q.Mode) + 1

	b := bytes.NewBuffer([]byte{})
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpRRQ) // write the operation, here it's a read
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(q.Filename) // write the Filename
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // write 0 byte
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode) // write mode
	if err != nil {
		return nil, err
	}
	err = b.WriteByte(0) // write 0 byte
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (q *ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code) // read opcode
	if err != nil {
		return err
	}

	if code != OpRRQ {
		return errors.New("invalid RRQ")
	}

	q.Filename, err = r.ReadString(0) // read filename, reads up to and including the 0 byte delimiter
	if err != nil {
		return err
	}

	q.Filename = strings.TrimRight(q.Filename, "\x00") // remove the 0 byte

	q.Mode, err = r.ReadString(0)
	if err != nil {
		return err
	}

	q.Mode = strings.TrimRight(q.Mode, "\x00") // remove the 0 byte

	if len(q.Mode) == 0 {
		return errors.New("invalid RRQ")
	}

	actual := strings.ToLower(q.Mode) // enforce octet mode
	if actual != "octet" {
		return errors.New("only binary transfers supported")
	}

	return nil
}

type Data struct {
	Block   uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	b.Grow(DatagramSize)

	// block numbers increment from 1
	d.Block++

	err := binary.Write(b, binary.BigEndian, OpData)
	if err != nil {
		return nil, err
	}
	err = binary.Write(b, binary.BigEndian, d.Block)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(b, d.Payload, BlockSize)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnmarshalBinary(p []byte) error {
	if l := len(p); l < HeaderSize || l > DatagramSize {
		return errors.New("invalid DATA")
	}

	var code OpCode

	// Read first two bytes from the payload, opcode size
	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpData {
		return errors.New("invalid data request")
	}

	// Read first next two bytes from the payload, block size
	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block)
	if err != nil {
		return err
	}

	// read the rest of the payload
	d.Payload = bytes.NewBuffer(p[4:])

	return nil
}

type Ack uint16

func (a Ack) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	buf.Grow(4) // 2 bytes for OpCode, 2 bytes for block number

	err := binary.Write(buf, binary.BigEndian, OpAck)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, a) // write the block number
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *Ack) UnmarshalBinary(p []byte) error {
	var op OpCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &op)
	if err != nil {
		return err
	}

	if op != OpAck {
		return errors.New("invalid ACK")
	}

	err = binary.Read(bytes.NewReader(p[2:]), binary.BigEndian, a)
	if err != nil {
		return err
	}

	return nil
}

type Err struct {
	Error   ErrCode
	Message string
}

func (e Err) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	buf.Grow(2 + 2 + len(e.Message) + 1) // 2 bytes for opcode, 2 bytes for errcode, message len, and 1 byte for null terminator

	err := binary.Write(buf, binary.BigEndian, OpErr)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, e.Error)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(e.Message)
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *Err) UnmarshalBinary(p []byte) error {
	var code OpCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpErr {
		return errors.New("invalid Err")
	}

	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &e.Error)
	if err != nil {
		return err
	}

	err = binary.Read(bytes.NewReader(p[4:]), binary.BigEndian, &e.Message)
	if err != nil {
		return err
	}

	e.Message = strings.TrimRight(e.Message, "\x00")

	return nil
}

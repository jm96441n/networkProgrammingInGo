package tftp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/charmbracelet/log"
)

type Server struct {
	Payload []byte        // payload served for all read request
	Retries uint8         // number of times to retry a failed  transaction
	Timeout time.Duration // the duration to wait for an  acknowledgement
}

type option func(*Server)

func NewServer(payload []byte, opts ...option) (Server, error) {
	if payload == nil {
		return Server{}, errors.New("payload is required")
	}
	s := Server{
		Payload: payload,
		Retries: 10,
		Timeout: 6 * time.Second,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s, nil
}

func WithRetries(r uint8) option {
	return func(s *Server) {
		s.Retries = r
	}
}

func WithTimeout(r time.Duration) option {
	return func(s *Server) {
		s.Timeout = r
	}
}

func (s Server) ListenAndServe(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	defer func() { _ = conn.Close() }()

	log.Info(fmt.Sprintf("Listening on %s...\n", conn.LocalAddr()))

	return s.Serve(conn)
}

func (s Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}

	var rrq ReadReq

	for {
		buf := make([]byte, DatagramSize)

		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}

		err = rrq.UnmarshalBinary(buf)

		if err != nil {
			log.Error(fmt.Sprintf("[%s] bad request: %v", addr, err))
			continue
		}
		go s.handle(addr.String(), rrq)
	}
}

func (s Server) handle(addr string, rrq ReadReq) {
	log.Info(fmt.Sprintf("[%s] requested file: %s", addr, rrq.Filename))

	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Error(fmt.Sprintf("[%s] dial: %v", addr, err))
		return
	}

	defer func() { _ = conn.Close() }()

	var (
		ackPkt  Ack
		errPkt  Err
		dataPkt = Data{Payload: bytes.NewReader(s.Payload)}
		buf     = make([]byte, DatagramSize)
	)

NEXTPACKET:
	for n := DatagramSize; n == DatagramSize; {

		data, err := dataPkt.MarshalBinary()
		if err != nil {
			log.Error(fmt.Sprintf("[%s] preparing data packet: %v", addr, err))
			return
		}
	RETRY:
		for i := s.Retries; i > 0; i-- {
			n, err = conn.Write(data) // send the packet
			if err != nil {
				log.Error(fmt.Sprintf("[%s] write: %v", addr, err))
				return
			}

			// wait for the client ack
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))

			_, err = conn.Read(buf)
			if err != nil {
				var netError net.Error
				// if we timeout then  retry
				if errors.As(err, &netError) && netError.Timeout() {
					continue RETRY
				}

				log.Error(fmt.Sprintf("[%s] wating for ACK: %v", addr, err))
				return
			}

			switch {
			case ackPkt.UnmarshalBinary(buf) == nil:
				if uint16(ackPkt) == dataPkt.Block {
					// received ack, send next packet
					continue NEXTPACKET
				}
			case errPkt.UnmarshalBinary(buf) == nil:
				log.Warn(fmt.Sprintf("[%s] received error: %v", addr, err))
				return
			default:
				log.Error(fmt.Sprintf("[%s] bad packet", addr))

			}
		}
		log.Error(fmt.Sprintf("[%s] exhausted retries", addr))
		return

	}
	log.Info(fmt.Sprintf("[%s] sent %d blocks", addr, dataPkt.Block))
}

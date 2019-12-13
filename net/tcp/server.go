package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/lj-team/go-generic/log"
)

type Handler interface {
	OnMessage([]byte) ([]byte, error)
}

type HandlerMaker func() Handler

func Server(addr string, maker HandlerMaker) error {
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	for {

		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		con := &connect{
			id:   <-nextId,
			addr: conn.RemoteAddr().String(),
			con:  conn,
		}

		proto := maker()

		go serverHandler(con, proto)
	}

	return nil
}

func serverHandler(conn *connect, proto Handler) {

	log.Trace(fmt.Sprintf("connection #%d (%s) opened", conn.id, conn.addr))

	defer conn.Close()

	br := bufio.NewReaderSize(conn, 10240)
	bw := bufio.NewWriterSize(conn, 10240)

	for {

		data, err := br.ReadBytes('\n')

		size := len(data)

		if size == 0 {
			break
		}

		data = bytes.TrimRight(data, "\r\n")

		if len(data) == 0 {
			continue
		}

		data, err = proto.OnMessage(data)
		if err != nil {
			break
		}

		if len(data) > 0 {

			_, err = bw.Write(data)
			if err != nil {
				break
			}
			if bw.WriteByte('\n') != nil {
				break
			}
			if bw.Flush() != nil {
				break
			}
		}

	}
}

package tcp

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/lj-team/go-generic/log"
)

var nextId chan int64

func init() {

	nextId = make(chan int64, 32)

	go func() {

		id := int64(1)

		for {
			nextId <- id
			id++
		}

	}()
}

type connect struct {
	id       int64
	addr     string
	con      net.Conn
	lastSend int64
}

func (c *connect) Close() {
	if c.con != nil {
		log.Trace(fmt.Sprintf("connection #%d (%s) closed", c.id, c.addr))
		c.con.Close()
		c.con = nil
	}
}

func (c *connect) Write(p []byte) (int, error) {

	now := time.Now().Unix()

	if c.lastSend+int64(KEEP_ALIVE) <= now {
		return 0, errors.New("keep alive")
	}

	c.lastSend = now

	c.con.SetDeadline(time.Now().Add(KEEP_ALIVE))
	return c.con.Write(p)
}

func (c *connect) Read(b []byte) (int, error) {
	c.con.SetDeadline(time.Now().Add(KEEP_ALIVE))
	return c.con.Read(b)
}

package udp

import (
	"net"
	"time"
)

type Client struct {
	con    *net.UDPConn
	addr   *net.UDPAddr
	buffer []byte
}

func NewClient(addr string) (*Client, error) {

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(setBufferSize)
	conn.SetWriteBuffer(setBufferSize)

	return &Client{con: conn, addr: raddr, buffer: make([]byte, 10240)}, nil
}

func (c *Client) Send(msg []byte) error {

	c.con.SetWriteDeadline(time.Now().Add(2 * time.Second))

	_, err := c.con.Write(msg)
	return err
}

func (c *Client) Read() ([]byte, error) {

	c.con.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, _, err := c.con.ReadFromUDP(c.buffer)
	if err != nil {
		return nil, err
	}

	return c.buffer[:n], nil
}

func (c *Client) Close() {
	c.con.Close()
}

package client

import (
	"bytes"
	"errors"
	"sync"

	"github.com/lj-team/go-generic/net/udp"
	"github.com/lj-team/go-generic/text/args"
)

type client struct {
	ua *udp.Client
	sync.Mutex
}

func New(addr string) (Client, error) {
	ua, err := udp.NewClient(addr)
	if err != nil {
		return nil, err
	}

	return &client{ua: ua}, nil
}

func (c *client) Close() {
	c.ua.Close()
}

func (c *client) send(msg []byte) error {
	c.Lock()
	defer c.Unlock()

	return c.ua.Send(msg)
}

func (c *client) sendread(msg []byte) ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	if err := c.ua.Send(msg); err != nil {
		return nil, err
	}

	data, err := c.ua.Read()
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(data, []byte("+")) {
		return data[1:], nil
	}

	if bytes.HasPrefix(data, []byte("-")) {
		return nil, errors.New(string(data[1:]))
	}

	return nil, errors.New("invalid answer message format")
}

func (c *client) Exec(command ...string) (string, error) {

	if len(command) == 0 {
		return "", errors.New("no command")
	}

	bb := bytes.NewBuffer(nil)

	bb.WriteString("exec")

	for _, v := range command {
		bb.WriteByte(' ')
		bb.WriteString(args.ToString(v))
	}

	data, err := c.sendread(bb.Bytes())
	return string(data), err
}

func (c *client) Async(command ...string) error {

	if len(command) == 0 {
		return errors.New("no command")
	}

	bb := bytes.NewBuffer(nil)

	bb.WriteString("async")

	for _, v := range command {
		bb.WriteByte(' ')
		bb.WriteString(args.ToString(v))
	}

	return c.send(bb.Bytes())
}

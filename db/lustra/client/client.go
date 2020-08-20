package client

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lj-team/go-generic/net/udp"
	"github.com/lj-team/go-generic/text/args"
)

type client struct {
	ua      *udp.Client
	startTS int64
	nextId  int64
	sync.Mutex
}

func New(addr string) (Client, error) {
	ua, err := udp.NewClient(addr)
	if err != nil {
		return nil, err
	}

	return &client{ua: ua, startTS: time.Now().Unix()}, nil
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

func (c *client) sendreadid(id string, msg []byte) ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	if err := c.ua.Send(msg); err != nil {
		return nil, err
	}

	for {

		data, err := c.ua.Read()
		if err != nil {
			return nil, err
		}

		positive := "+" + id

		if bytes.HasPrefix(data, []byte(positive)) {
			data = data[len(positive):]
			return data, nil
		}

		negative := "-" + id

		if bytes.HasPrefix(data, []byte(negative)) {
			return nil, errors.New(string(data[len(negative):]))
		}
	}

	return nil, errors.New("invalid answer message format")
}

func (c *client) Exec(command ...string) (string, error) {

	if len(command) == 0 {
		return "", errors.New("no command")
	}

	atomic.AddInt64(&c.nextId, 1)

	id := fmt.Sprintf("%04x%04x", c.startTS&0xffff, c.nextId&0xffff)

	bb := bytes.NewBuffer(nil)

	bb.WriteString("exec" + id)

	for _, v := range command {
		bb.WriteByte(' ')
		bb.WriteString(args.ToString(v))
	}

	data, err := c.sendreadid(id, bb.Bytes())
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

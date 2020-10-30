package proxy

import (
	"bytes"
	"errors"

	"github.com/lj-team/go-generic/db/lustra/client"
	"github.com/lj-team/go-generic/text/args"

	_ "github.com/lj-team/go-generic/db/lustra/storage/engine/cache"
)

type Proxy struct {
	cons []client.Client
}

func New(addrs []string) *Proxy {

	if len(addrs) < 1 {
		return nil
	}

	p := &Proxy{}

	for _, addr := range addrs {
		ua, err := client.New(addr)
		if err != nil {
			return nil
		}

		p.cons = append(p.cons, ua)
	}

	return p
}

func (px *Proxy) Close() {
	for _, c := range px.cons {
		c.Close()
	}
}

func NewStub() *Proxy {

	con, _ := client.NewStub()

	p := &Proxy{
		cons: []client.Client{con},
	}

	return p
}

func (p *Proxy) Exec(command ...string) (string, error) {

	var res string
	var err error = errors.New("no connects")

	has := false

	for _, con := range p.cons {

		if !has {
			res, err = con.Exec(command...)
			if err == nil {
				has = true
			}
		} else {
			con.Exec(command...)
		}

	}

	return res, err
}

func (p *Proxy) Fetch(command ...string) (string, error) {
	var res string
	var err error = errors.New("no connects")

	for _, con := range p.cons {

		res, err = con.Exec(command...)
		if err == nil {
			return res, nil
		}
	}

	return res, err
}

func (p *Proxy) Async(command ...string) {

	for _, con := range p.cons {
		con.Async(command...)
	}

}

func (p *Proxy) CommandString(command ...string) string {

	bb := bytes.NewBuffer(nil)

	for i, v := range command {
		if i > 0 {
			bb.WriteByte(' ')
		}
		bb.WriteString(args.ToString(v))
	}

	return bb.String()
}

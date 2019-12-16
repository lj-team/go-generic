package client

import (
	"bytes"
	"errors"
	"sync"

	"github.com/lj-team/go-generic/db/lustra/server/handler"
	"github.com/lj-team/go-generic/db/lustra/storage"
	"github.com/lj-team/go-generic/text/args"

	_ "github.com/lj-team/go-generic/db/lustra/storage/engine/cache"
)

type stub struct {
	sync.Mutex
	dbh *storage.Storage
}

func NewStub() (Client, error) {

	dbh, err := storage.Open("cache", "size=10240 nodes=4")
	if err != nil {
		return nil, err
	}

	return &stub{dbh: dbh}, nil
}

func (s *stub) Close() {
	s.dbh.Close()
}

func (s *stub) proc(msg []byte) ([]byte, error) {

	msg = handler.Handler(s.dbh, msg)

	if len(msg) > 0 {

		if bytes.HasPrefix(msg, []byte("+")) {
			return msg[1:], nil
		}

		if bytes.HasPrefix(msg, []byte("-")) {
			return nil, errors.New(string(msg[1:]))
		}

		return msg, errors.New("invalid answer message format")
	}

	return []byte("ok"), nil
}

func (s *stub) Exec(command ...string) (string, error) {

	s.Lock()
	defer s.Unlock()

	if len(command) == 0 {
		return "", errors.New("no command")
	}

	bb := bytes.NewBuffer(nil)

	bb.WriteString("exec")

	for _, v := range command {
		bb.WriteByte(' ')
		bb.WriteString(args.ToString(v))
	}

	data, err := s.proc(bb.Bytes())
	return string(data), err
}

func (s *stub) Async(command ...string) error {

	s.Lock()
	defer s.Unlock()

	if len(command) == 0 {
		return errors.New("no command")
	}

	bb := bytes.NewBuffer(nil)

	bb.WriteString("async")

	for _, v := range command {
		bb.WriteByte(' ')
		bb.WriteString(args.ToString(v))
	}

	s.proc(bb.Bytes())

	return nil
}

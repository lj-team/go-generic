package client

import (
	"testing"

	"github.com/lj-team/go-generic/db/lustra/global"
)

func TestStub(t *testing.T) {

	ua, err := NewStub()
	if err != nil {
		t.Fatal("NewStub failed")
	}
	defer ua.Close()

	res, e := ua.Exec("version")
	if e != nil || res != global.Version {
		t.Fatal("Exec version failed")
	}

	if _, e := ua.Exec(); e == nil {
		t.Fatal("Invalid Exec call not return error")
	}

	e = ua.Async("set", "1", "2")
	if e != nil {
		t.Fatal("Async set failed")
	}

	if ua.Async() == nil {
		t.Fatal("Invalid Async call not return error")
	}

	if _, e := ua.Exec("get", "1", "2", "3"); e == nil {
		t.Fatal("Invalid Exec get wrong params call not return error")
	}

	res, e = ua.Exec("get", "1")
	if e != nil || res != "2" {
		t.Fatal("Exec get failed")
	}
}

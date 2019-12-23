package server

import (
	"testing"
	"time"

	"github.com/lj-team/go-generic/db/lustra/global"
	"github.com/lj-team/go-generic/db/lustra/proxy"

	_ "github.com/lj-team/go-generic/db/lustra/storage/engine/cache"
)

func TestServer(t *testing.T) {

	go func() {
		err := Start("127.0.0.1:10002", "cache", "size=1024 nodes=4")
		if err != nil {
			panic(err)
		}
	}()

	<-time.After(time.Millisecond * 10)

	px := proxy.New([]string{"127.0.0.1:10002"})

	tF := func(exec bool, cmd []string, answer string, err string, wait time.Duration) {

		if exec {
			str, e := px.Exec(cmd...)
			if str != answer {
				t.Fatal("invalid answer cmd=", cmd, " ret=", str, " wait=", answer)
			}

			if e != nil {
				if err == "" {
					t.Fatal("Not expected error cmd=", cmd)
				}
			} else {
				if err != "" {
					t.Fatal("Expect error cmd=", cmd, " err=", err)
				}
			}
		} else {
			px.Async(cmd...)
			if wait != 0 {
				<-time.After(wait)
			}
		}

	}

	tF(true, []string{"ping"}, "pong", "", 0)
	tF(true, []string{"ping", "1"}, "", "invalid params", 0)
	tF(false, []string{"ping"}, "pong", "", 0)
	tF(true, []string{"version"}, global.Version, "", 0)
	tF(false, []string{"version", "1"}, "", "invalid params", 0)
	tF(true, []string{"set", "1", "01"}, "ok", "", 0)
	tF(true, []string{"get", "1"}, "01", "", 0)
	tF(false, []string{"set", "1", "02", "2", "03"}, "ok", "", time.Millisecond*300)
	tF(true, []string{"get", "1"}, "02", "", 0)
	tF(true, []string{"get", "2"}, "03", "", 0)
	tF(true, []string{"del", "1"}, "ok", "", 0)
	tF(true, []string{"get", "1"}, "", "", 0)
	tF(true, []string{"inc", "1"}, "1", "", 0)
	tF(true, []string{"inc", "1"}, "2", "", 0)
	tF(true, []string{"inc", "1"}, "3", "", 0)
	tF(true, []string{"dec", "1"}, "2", "", 0)
	tF(true, []string{"dec", "1"}, "1", "", 0)
	tF(true, []string{"dec", "1"}, "0", "", 0)
	tF(true, []string{"dec", "1"}, "0", "", 0)
	tF(true, []string{"setnx", "100", "1"}, "1", "", 0)
	tF(true, []string{"setnx", "100", "2"}, "1", "", 0)
	tF(true, []string{"setifmore", "100", "0"}, "1", "", 0)
	tF(true, []string{"setifmore", "100", "1"}, "1", "", 0)
	tF(true, []string{"setifmore", "100", "2"}, "2", "", 0)
	tF(true, []string{"setifmore", "100", "1"}, "2", "", 0)
	tF(true, []string{"setifless", "100", "1"}, "1", "", 0)
	tF(true, []string{"setifless", "100", "0"}, "", "", 0)
	tF(true, []string{"setnx", "100"}, "", "invalid params", 0)
	tF(true, []string{"incby", "1", "12"}, "12", "", 0)
	tF(true, []string{"incby", "1", "12"}, "24", "", 0)
	tF(true, []string{"incby", "1", "-12"}, "", "invalid params", 0)
	tF(true, []string{"decby", "1", "18"}, "6", "", 0)
	tF(true, []string{"decby", "1", "18"}, "0", "", 0)
	tF(true, []string{"decby", "1", "18"}, "0", "", 0)
	tF(true, []string{"decby", "1", "-18"}, "", "invalid params", 0)
	tF(true, []string{"cbadd", "1", "1", "5"}, "ok", "", 0)
	tF(true, []string{"cbadd", "1", "2", "5"}, "ok", "", 0)
	tF(true, []string{"cbadd", "1", "3", "5"}, "ok", "", 0)
	tF(true, []string{"cbadd", "1", "4", "5"}, "ok", "", 0)
	tF(true, []string{"cbadd", "1", "5", "5"}, "ok", "", 0)
	tF(true, []string{"cbadd", "1", "6", "5"}, "ok", "", 0)
	tF(true, []string{"get", "1"}, `["2","3","4","5","6"]`, "", 0)
	tF(true, []string{"hset", "h1", "1", "2"}, `ok`, "", 0)
	tF(true, []string{"hset", "h1", "3", "5"}, `ok`, "", 0)
	tF(true, []string{"hset", "h1", "5", "8"}, `ok`, "", 0)
	tF(true, []string{"hset", "h1", "8", "13"}, `ok`, "", 0)
	tF(true, []string{"hset", "h1", "8", "13", "."}, ``, "invalid params", 0)
	tF(true, []string{"hget", "h1", "8"}, `13`, "", 0)
	tF(true, []string{"hdel", "h1", "8"}, `ok`, "", 0)
	tF(true, []string{"hget", "h1", "8"}, ``, "", 0)
	tF(true, []string{"hinc", "h1", "8"}, `1`, "", 0)
	tF(true, []string{"hinc", "h1", "8", "0"}, ``, "invalid params", 0)
	tF(true, []string{"hinc", "h1", "8"}, `2`, "", 0)
	tF(true, []string{"hinc", "h1", "8"}, `3`, "", 0)
	tF(true, []string{"hdec", "h1", "8"}, `2`, "", 0)
	tF(true, []string{"hdec", "h1", "8"}, `1`, "", 0)
	tF(true, []string{"hdec", "h1", "8"}, `0`, "", 0)
	tF(true, []string{"hdec", "h1", "8", "0"}, ``, "invalid params", 0)
	tF(true, []string{"hdec", "h1", "8"}, `0`, "", 0)
	tF(true, []string{"hget", "h1", "8"}, ``, "", 0)
	tF(true, []string{"hget", "h1", "8", "1"}, ``, "invalid params", 0)
	tF(true, []string{"hincby", "h1", "5", "-10"}, ``, "invalid params", 0)
	tF(true, []string{"hincby", "h1", "5", "10"}, `18`, "", 0)
	tF(true, []string{"hdecby", "h1", "5", "10"}, `8`, "", 0)
	tF(true, []string{"hdecby", "h1", "5", "-10"}, ``, "invalid params", 0)
	tF(true, []string{"hdecby", "h1", "5", "10"}, `0`, "", 0)
	tF(true, []string{"hsetnx", "h1", "5", "10"}, `10`, "", 0)
	tF(true, []string{"hsetnx", "h1", "5", "11"}, `10`, "", 0)
	tF(true, []string{"hsetnx", "h1", "5", "11", "19"}, ``, "invalid params", 0)
	tF(true, []string{"hsetifmore", "h1", "5", "8"}, "10", "", 0)
	tF(true, []string{"hsetifmore", "h1", "5", "13"}, "13", "", 0)
	tF(true, []string{"hsetifless", "h1", "5", "23"}, "13", "", 0)
	tF(true, []string{"hsetifless", "h1", "5", "3"}, "3", "", 0)
	tF(true, []string{"hsetifless", "h1", "5", "0"}, "", "", 0)
}

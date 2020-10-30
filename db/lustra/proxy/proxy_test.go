package proxy

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lj-team/go-generic/db/lustra/global"
)

func TestProxy(t *testing.T) {

	px := NewStub()

	tE := func(exec bool, cmd []string, answer string, err string, wait time.Duration) {

		if exec {
			str, e := px.Exec(cmd...)
			if str != answer {
				fmt.Println(str)
				t.Fatal("invalid answer cmd=", cmd, " wait=", answer)
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

	tF := func(exec bool, cmd []string, answer string, err string, wait time.Duration) {

		if exec {
			str, e := px.Fetch(cmd...)
			if str != answer {
				fmt.Println(str)
				t.Fatal("invalid answer cmd=", cmd, " wait=", answer)
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

	tB := func(cmd []string, res string) {
		if px.CommandString(cmd...) != res {
			t.Fatalf("Batch failed for: %s", strings.Join(cmd, " "))
		}
	}

	tE(true, []string{"ping"}, "pong", "", 0)
	tE(true, []string{"ping", "1"}, "", "invalid params", 0)
	tE(false, []string{"ping"}, "pong", "", 0)

	tF(true, []string{"version"}, global.Version, "", 0)
	tF(false, []string{"version", "1"}, "", "invalid params", 0)

	tB([]string{"version"}, "version")
	tB([]string{"get", "123"}, "get 123")
	tB([]string{"hset", "1.2", "code", ""}, "hset 1.2 code \"\"")
}

package lustra

import (
	"strconv"
	"testing"
	"time"

	"github.com/lj-team/go-generic/db/lustra/proxy"
)

func TestLustar(t *testing.T) {

	go func() {

		err := Start("127.0.0.1:10001", "cache", "size=1024 nodes=4")
		if err != nil {
			panic(err)
		}
	}()

	<-time.After(time.Millisecond * 10)

	px := proxy.New([]string{"127.0.0.1:10001"})

	tS := func(k, w string) {
		px.Exec("set", k, w)

		res, _ := px.Exec("get", k)

		if res != w {
			t.Fatal("Set not work for: " + k)
		}
	}

	for i := 0; i < 1000; i++ {
		tS(strconv.Itoa(i), strconv.Itoa(i*10))
	}

	tI := func(k string, inc int64, wait int64) {

		var res string

		if inc > 0 {
			res, _ = px.Exec("incby", k, strconv.FormatInt(inc, 10))
		} else {
			res, _ = px.Exec("decby", k, strconv.FormatInt(-inc, 10))
		}

		if res != strconv.FormatInt(wait, 10) {
			t.Fatal("Invalid Inc: " + k)
		}
	}

	tI("_", 1, 1)
	tI("_", 2, 3)
	tI("_", 3, 6)
	tI("_", -4, 2)
	tI("_", -3, 0)
	tI("_", 3, 3)
}

package counter

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/lj-team/go-generic/time/strftime"
)

type Counter struct {
	val  int64
	mod  int64
	file string
}

func New(filename string, mod int64) *Counter {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte("0")
	}

	if mod <= 0 {
		mod = math.MaxInt64
	}

	c := &Counter{file: filename, val: 0, mod: mod}

	str := strings.TrimSpace(string(data))

	c.val, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		c.val = 0
	}

	if c.val < 0 {
		c.val = 0
	}

	c.val = c.val % c.mod

	return c

}

func (c *Counter) flush() {
	data := strconv.FormatInt(c.val, 10)
	ioutil.WriteFile(c.file, []byte(data), 0644)
}

func (c *Counter) Set(val int64) {

	if val < 0 {
		val = 0
	}

	c.val = val % c.mod
	c.flush()
}

func (c *Counter) Inc() int64 {
	c.val = (c.val + 1) % c.mod
	c.flush()
	return c.val
}

func (c *Counter) SetDate(t time.Time) {
	str := strftime.Format("%Y%m%d", t)
	val, _ := strconv.ParseInt(str, 10, 64)
	c.Set(val)
}

func (c *Counter) Get() int64 {
	return c.val
}

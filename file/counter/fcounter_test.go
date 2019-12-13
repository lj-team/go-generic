package counter

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/lj-team/go-generic/time/strftime"
)

func TestCounter(t *testing.T) {

	content := []byte("0")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	cnt := New(tmpfile.Name(), 10)

	if cnt.Get() != 0 {
		t.Fatal("Init Get error")
	}

	for i := int64(0); i < 25; i++ {
		if cnt.Inc() != (i+1)%10 {
			t.Fatal("Inc not work")
		}
	}

	data, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err.Error())
	}

	if bytes.Compare(data, []byte("5")) != 0 {
		t.Fatal("invalid file content")
	}

	cnt.Set(18)

	data, err = ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err.Error())
	}

	if bytes.Compare(data, []byte("8")) != 0 {
		t.Fatal("invalid file content")
	}

	if cnt.Get() != 8 {
		t.Fatal("invalid get value")
	}

	tm := time.Now()

	cnt2 := New(tmpfile.Name(), -1)

	if cnt2.Get() != 8 {
		t.Fatal("invalid loaded value")
	}

	cnt2.SetDate(tm)

	val, _ := strconv.ParseInt(strftime.Format("%Y%m%d", tm), 10, 64)

	if cnt2.Get() != val {
		t.Fatal("Invalid counter value")
	}

	cnt2.Set(-1)

	if cnt2.Get() != 0 {
		t.Fatal("invalid value")
	}
}

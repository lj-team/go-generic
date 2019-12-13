package cache

import (
	"testing"
	"time"
)

func TestItem(t *testing.T) {

	i1 := item{expire: time.Now().Unix() - 1}

	if i1.isAlive() {
		t.Fatal("isAlive failed")
	}

	i2 := item{expire: time.Now().Unix() + 5}

	if !i2.isAlive() {
		t.Fatal("isAlive failed")
	}
}

package mapcnt

import (
	"testing"
)

func TestFromMapCnt(t *testing.T) {

	src := map[string]float64{"1": 12.1, "2": 11.2, "3": 13.3, "4": 14.4}

	res := List(src, &ListOpts{MinVal: 12, Limit: 10}).([]string)

	wait := []string{"4", "3", "1"}

	if len(wait) != len(res) {
		t.Fatal("List wrong result size")
	}

	for i, v := range res {
		if wait[i] != v {
			t.Fatal("List wron result item")
		}
	}

	res = List(src, nil).([]string)

	wait = []string{"4", "3", "1", "2"}

	if len(wait) != len(res) {
		t.Fatal("slice.FromMapCnt wrong result size")
	}

	for i, v := range res {
		if wait[i] != v {
			t.Fatal("slice.FromMapCnt wron result item")
		}
	}
}

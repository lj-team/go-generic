package pack

import (
	"testing"
)

func TestInt2Bytes(t *testing.T) {
	val := int64(1234567890)
	b := Int2Bytes(val)
	if val != Bytes2Int(b) {
		t.Fatal(Bytes2Int(b))
		t.Fatal("expecting 1234567890")
	}

}

func TestIntList2Bytes(t *testing.T) {
	list := []int64{1, 129, 1045, 76532, 462784628}
	b := IntList2Bytes(list)
	res := Bytes2IntList(b)

	if len(res) != len(list) {
		t.Fatal("ivalid decoded size")
	}

	for i, cur := range res {
		if cur != list[i] {
			t.Fatalf("expected %d", cur)
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	v1 := int64(64)
	v2 := int32(32)
	v3 := int16(16)
	v4 := float32(32.00)
	v5 := float64(64.00)

	data := Encode(v1, v2, v3, v4, v5)

	if len(data) != 26 {
		t.Fatalf("wrong bytes number")
	}

	v1 = int64(0)
	v2 = int32(0)
	v3 = int16(0)
	v4 = float32(0)
	v5 = float64(0)

	if Decode(data, &v1, &v2, &v3, &v4, &v5) != nil {
		t.Fatal("decode error")
	}

	if v1 != int64(64) || v2 != int32(32) || v3 != int16(16) || v4 != float32(32.00) || v5 != float64(64.00) {
		t.Fatal("decode errors")
	}
}

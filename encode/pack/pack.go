package pack

import (
	"bytes"
	"encoding/binary"
)

func Encode(list ...interface{}) []byte {
	b := new(bytes.Buffer)
	for _, item := range list {

		switch v := item.(type) {
		case string:
			binary.Write(b, binary.BigEndian, []byte(v))
		default:
			binary.Write(b, binary.BigEndian, item)
		}
	}
	return b.Bytes()
}

func Decode(data []byte, list ...interface{}) error {
	reader := bytes.NewReader(data)
	for _, item := range list {
		if err := binary.Read(reader, binary.BigEndian, item); err != nil {
			return err
		}
	}

	return nil
}

func Int2Bytes(val int64) []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, val)
	return b.Bytes()
}

func IntList2Bytes(vals []int64) []byte {
	if vals == nil {
		return make([]byte, 0)
	}

	b := new(bytes.Buffer)

	for _, val := range vals {
		binary.Write(b, binary.BigEndian, val)
	}

	return b.Bytes()
}

func Bytes2Int(b []byte) int64 {
	var val int64
	reader := bytes.NewReader(b)
	err := binary.Read(reader, binary.BigEndian, &val)
	if err != nil {
		val = 0
	}
	return val
}

func Bytes2IntList(b []byte) []int64 {
	if b == nil {
		return make([]int64, 0)
	}

	var val int64

	size := int(len(b) / 8)
	res := make([]int64, size)
	reader := bytes.NewReader(b)

	for i := 0; i < size; i++ {
		binary.Read(reader, binary.BigEndian, &val)
		res[i] = val
	}

	return res
}

package chash

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2017-11-01

import (
	"hash/crc32"
)

type Hash struct {
	num  int
	step uint32
}

func New(nodes int) *Hash {
	h := &Hash{
		num:  nodes,
		step: uint32(0xffffffff) / uint32(nodes),
	}

	return h
}

func (h *Hash) Get(key []byte) int {

	if h.num < 2 {
		return 0
	}

	v := crc32.ChecksumIEEE(key)
	num := int(v / h.step)
	if num >= h.num {
		num = h.num - 1
	}

	return num
}

func (h *Hash) Next(num int) int {

	if h.num < 2 {
		return 0
	}

	return (num + 1) % h.num
}

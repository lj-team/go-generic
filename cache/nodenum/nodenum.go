package nodenum

import (
	"hash/crc32"
)

type NodeNum struct {
	num int
}

func New(nodes int) *NodeNum {
	h := &NodeNum{
		num: nodes,
	}

	return h
}

func (h *NodeNum) Get(key string) int {

	if h.num < 2 {
		return 0
	}

	v := crc32.ChecksumIEEE([]byte(key))

	return int(v % uint32(h.num))
}

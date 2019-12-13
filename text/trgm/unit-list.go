package trgm

type UintList []uint64

func (l UintList) Len() int {
	return len(l)
}

func (l UintList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l UintList) Less(i, j int) bool {
	return l[i] < l[j]
}

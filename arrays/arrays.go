package arrays

func Filter[T comparable](arr []T, cond func(T) bool) []T {
	result := make([]T, 0, len(arr))
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}

func In[T comparable](value T, arr []T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func ChunkBy[T any](items []T, chunkSize int) [][]T {
	var _chunks = make([][]T, 0, (len(items)/chunkSize)+1)
	for chunkSize < len(items) {
		items, _chunks = items[chunkSize:], append(_chunks, items[0:chunkSize:chunkSize])
	}
	return append(_chunks, items)
}

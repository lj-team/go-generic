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
	for i := range arr {
		if arr[i] == value {
			return true
		}
	}
	return false
}

func OneIn[T comparable](values []T, arr []T) bool {
	for i := range arr {
		for j := range values {
			if arr[i] == values[j] {
				return true
			}
		}
	}
	return false
}

func ChunkBy[T any](items []T, chunkSize int) [][]T {
	var chunks = make([][]T, 0, (len(items)/chunkSize)+1)
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

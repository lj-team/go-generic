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

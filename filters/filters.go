package filters

func String(arr []string, cond func(string) bool) []string {
	result := make([]string, 0, len(arr))
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}

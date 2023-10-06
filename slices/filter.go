package slices

func Filter[T any](s []T, filterFunc func(int, T) bool) []T {
	var out []T
	for i, v := range s {
		if filterFunc(i, v) {
			out = append(out, v)
		}
	}
	return out
}

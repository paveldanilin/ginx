package slices

func First[T any](s []T, matchFunc func(T) bool) (T, bool) {
	for _, v := range s {
		if matchFunc(v) {
			return v, true
		}
	}
	var empty T
	return empty, false
}

package slices

func Map[T, U any](s []T, f func(T) U) []U {
	var out []U

	for _, e := range s {
		out = append(out, f(e))
	}

	return out
}

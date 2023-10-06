package slices

func Join[T any](a []T, b []T) []T {
	var out []T

	for _, v := range a {
		out = append(out, v)
	}

	for _, v := range b {
		out = append(out, v)
	}

	return out
}

package transform

func Map[T, V any](elems []T, fn func(T) V) []V {
	result := make([]V, len(elems))
	for i, t := range elems {
		result[i] = fn(t)
	}
	return result
}

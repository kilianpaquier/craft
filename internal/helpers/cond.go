package helpers

// Or is a simple If Then Else function allowing inlining simple conditions.
func Or[T any](cond bool, a, b T) T { //nolint:revive
	if cond {
		return a
	}
	return b
}

// Coalesce returns the first non zero value from inputs.
func Coalesce[T comparable](values ...T) T {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}
	return zero
}

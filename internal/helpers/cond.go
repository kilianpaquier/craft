package helpers

// Or is a simple If Then Else function allowing inlining simple conditions.
func Or[T any](cond bool, a, b T) T { //nolint:revive
	if cond {
		return a
	}
	return b
}

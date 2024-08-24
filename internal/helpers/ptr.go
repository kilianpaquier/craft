package helpers

// ToPtr takes an input value and returns its associated ptr.
//
// Note that as the input may not be a ptr, a copy may be done.
func ToPtr[T any](in T) *T {
	return &in
}

// FromPtr takes an input ptr to a value and returns its value.
//
// In case the ptr is nil, a zero value is returned.
func FromPtr[T any](in *T) (out T) {
	if in == nil {
		return out
	}
	return *in
}

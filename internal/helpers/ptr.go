package helpers

// ToPtr takes an input value and returns its associated ptr.
//
// Note that as the input may not be a ptr, a copy may be done.
func ToPtr[T any](v T) *T {
	return &v
}

// FromPtr takes an input ptr to a value and returns its value.
//
// In case the ptr is nil, a zero value is returned.
func FromPtr[T any](v *T) (zero T) {
	if v == nil {
		return zero
	}
	return *v
}

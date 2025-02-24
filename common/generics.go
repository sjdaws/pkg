package common

// Pointer convert input to pointer.
func Pointer[T any](value T) *T {
	return &value
}

// True check whether a value is true.
func True(value *bool) bool {
	// Nil is false
	return !(value == nil || *value == false) //nolint:gosimple
}

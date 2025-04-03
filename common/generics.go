package common

// Pointer convert input to pointer.
func Pointer[T any](value T) *T {
	return &value
}

// True check whether a value is true.
func True(value *bool) bool {
	//nolint:staticcheck // Nil is false, but we need to check it exists before accessing pointer
	return !(value == nil || *value == false)
}

package common

// True check whether a value is true.
func True(value *bool) bool {
	//nolint:staticcheck // Nil is false, but we need to check it exists before accessing pointer
	return !(value == nil || *value == false)
}

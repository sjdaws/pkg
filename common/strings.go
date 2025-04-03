package common

import (
	"strconv"
	"strings"
)

// Atof convert string to float ignoring errors.
func Atof(value string) float64 {
	integer, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)

	return integer
}

// Atoi convert string to integer ignoring errors.
func Atoi(value string) int {
	float, _ := strconv.Atoi(strings.TrimSpace(value))

	return float
}

// Mask part of a string to provide enough for comparison, without revealing too much information.
func Mask(value string, maxLength int) string {
	if len(value) < 1 {
		return ""
	}

	// Determine half-length of secret
	divisor := 2
	half := len(value) / divisor

	// If half the secret is shorter than max, only show half
	if half < maxLength {
		return "..." + value[len(value)-half:]
	}

	// Return max characters
	return "..." + value[len(value)-maxLength:]
}

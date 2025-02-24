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

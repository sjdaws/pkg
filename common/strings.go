package common

import (
	"strconv"
	"strings"

	"github.com/fatih/camelcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/sjdaws/pkg/uuid"
)

// Atof convert string to float ignoring errors.
func Atof(value string) float64 {
	float, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)

	return float
}

// Atoi convert string to integer ignoring errors.
func Atoi(value string) int {
	integer, _ := strconv.Atoi(strings.TrimSpace(value))

	return integer
}

// Atou convert string to uuid ignoring errors.
func Atou(value string) uuid.UUID {
	parsed, _ := uuid.Parse(strings.TrimSpace(value))

	return parsed
}

// Atou64 convert string to uint64 ignoring errors.
func Atou64(value string) uint64 {
	integer, _ := strconv.ParseUint(value, 10, 64)

	return integer
}

// FriendlyName converts CamelCaseStructFields into Camel case struct fields.
func FriendlyName(name string) string {
	// Break string based on camel case
	words := camelcase.Split(name)

	// Process each word, and lowercase it unless it's all uppercase
	for index, word := range words {
		if word == strings.ToUpper(word) {
			continue
		}

		// Convert word to all lower case
		words[index] = strings.ToLower(word)

		// Title case the first word only
		if index == 0 {
			words[index] = cases.Title(language.English).String(words[index])
		}
	}

	return strings.Join(words, " ")
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

// Options take space separated words, add quotes around each word and returns a comma separated string. The final
// option can be prefixed with an inclusive/exclusive prefix, e.g. or/and.
func Options(options string, finalPrefix string) string {
	values := strings.Split(options, " ")
	words := make([]string, 0)

	for _, value := range values {
		words = append(words, "'"+value+"'")
	}

	if finalPrefix != "" && len(words) > 1 {
		words[len(words)-1] = finalPrefix + " " + words[len(words)-1]
	}

	return strings.Join(words, ", ")
}

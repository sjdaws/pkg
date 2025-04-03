package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/common"
)

func TestAtoi(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected int
		value    string
	}{
		"at start": {
			expected: 0,
			value:    "8string starts with a number",
		},
		"contains numbers": {
			expected: 0,
			value:    "string contains1 some numb2ers",
		},
		"invalid": {
			expected: 0,
			value:    "string contains no numbers",
		},
		"negative float": {
			expected: -0,
			value:    "-5.24",
		},
		"negative integer": {
			expected: -4,
			value:    "-4",
		},
		"valid float": {
			expected: 0,
			value:    "3.290",
		},
		"valid number": {
			expected: 1,
			value:    "1",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := common.Atoi(testcase.value)
			assert.Equal(t, testcase.expected, actual)
		})
	}
}

func TestAtof(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected float64
		value    string
	}{
		"at start": {
			expected: 0,
			value:    "8string starts with a number",
		},
		"contains numbers": {
			expected: 0,
			value:    "string contains1 some numb2ers",
		},
		"invalid": {
			expected: 0,
			value:    "string contains no numbers",
		},
		"negative float": {
			expected: -5.24,
			value:    "-5.24",
		},
		"negative integer": {
			expected: -4,
			value:    "-4",
		},
		"valid float": {
			expected: 3.29,
			value:    "3.290",
		},
		"valid number": {
			expected: 1,
			value:    "1",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := common.Atof(testcase.value)

			if testcase.expected == 0 {
				assert.Zero(t, actual)

				return
			}

			assert.InEpsilon(t, testcase.expected, actual, 0.0001)
		})
	}
}

func TestMask(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected  string
		maxLength int
		secret    string
	}{
		"no secret": {
			expected:  "",
			maxLength: 6,
			secret:    "",
		},
		"short secret": {
			expected:  "...rt",
			maxLength: 6,
			secret:    "short",
		},
		"long secret": {
			expected:  "...secret",
			maxLength: 6,
			secret:    "a really long secret",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := common.Mask(testcase.secret, testcase.maxLength)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}

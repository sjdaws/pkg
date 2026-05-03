package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sjdaws.com/pkg/common"
)

func TestTrue(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected bool
		value    *bool
	}{
		"false": {
			expected: false,
			value:    new(false),
		},
		"nil": {
			expected: false,
			value:    nil,
		},
		"true": {
			expected: true,
			value:    new(true),
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testcase.expected, common.True(testcase.value))
		})
	}
}

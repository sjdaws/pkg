package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/pkg/common"
)

func TestPointer(t *testing.T) {
	t.Parallel()

	actual := common.Pointer(true)

	value := true
	expected := &value

	assert.Equal(t, expected, actual)
}

func TestTrue(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected bool
		value    *bool
	}{
		"false": {
			expected: false,
			value:    common.Pointer(false),
		},
		"nil": {
			expected: false,
			value:    nil,
		},
		"true": {
			expected: true,
			value:    common.Pointer(true),
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testcase.expected, common.True(testcase.value))
		})
	}
}

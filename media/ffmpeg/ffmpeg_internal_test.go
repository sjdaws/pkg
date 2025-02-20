package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForce(t *testing.T) {
	t.Parallel()

	conversion := Input("").Output("")

	require.False(t, conversion.force)

	conversion.Force()

	assert.True(t, conversion.force)
}

func Test_getCommandArguments(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		error          string
		expected       []string
		force          bool
		inputFilename  string
		inputOptions   []Option
		outputFilename string
		outputOptions  []Option
	}{
		"forced": {
			expected:       []string{"-hide_banner", "-v", "error", "-i", "input.mp4", "output.mp4", "-progress", "/dev/stdout", "-y"},
			force:          true,
			inputFilename:  "input.mp4",
			outputFilename: "output.mp4",
		},
		"no input filename": {
			error: "input filename is required",
		},
		"no output filename": {
			error:         "output filename is required",
			inputFilename: "input.mp4",
		},
		"no options": {
			expected:       []string{"-hide_banner", "-v", "error", "-i", "input.mp4", "output.mp4", "-progress", "/dev/stdout"},
			inputFilename:  "input.mp4",
			outputFilename: "output.mp4",
		},
		"options": {
			expected:       []string{"-hide_banner", "-v", "error", "-ignore_unknown", "-ss", "00:00:00", "-i", "input.mp4", "-crf", "23", "output.mp4", "-progress", "/dev/stdout", "-y"},
			force:          true,
			inputFilename:  "input.mp4",
			inputOptions:   []Option{{Key: "ignore_unknown", Value: nil}, {Key: "ss", Value: "00:00:00"}},
			outputFilename: "output.mp4",
			outputOptions:  []Option{{Key: "crf", Value: "23"}},
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			conversion := Input(testcase.inputFilename, testcase.inputOptions...).
				Output(testcase.outputFilename, testcase.outputOptions...)

			if testcase.force {
				conversion.Force()
			}

			actual, err := conversion.getCommandArguments()

			if testcase.error != "" {
				require.Error(t, err)

				require.EqualError(t, err, testcase.error)
				assert.Nil(t, actual)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, testcase.expected, actual)
		})
	}
}

func Test_processOptions(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected []string
		options  []Option
	}{
		"invalid value - empty": {
			expected: []string{"-key"},
			options:  []Option{{Key: "key", Value: ""}},
		},
		"invalid value - quoted whitespace": {
			expected: []string{"-key"},
			options:  []Option{{Key: "key", Value: `" " `}},
		},
		"invalid value - quotes": {
			expected: []string{"-key"},
			options:  []Option{{Key: "key", Value: `""`}},
		},
		"invalid value - whitespace": {
			expected: []string{"-key"},
			options:  []Option{{Key: "key", Value: " "}},
		},
		"key and value": {
			expected: []string{"-key", "value"},
			options:  []Option{{Key: "key", Value: "value"}},
		},
		"key only": {
			expected: []string{"-key"},
			options:  []Option{{Key: "key", Value: nil}},
		},
		"no key": {
			expected: []string{},
			options:  []Option{{Key: "", Value: "value"}},
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			conversion := Input("").Output("")

			actual := conversion.processOptions(testcase.options)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}

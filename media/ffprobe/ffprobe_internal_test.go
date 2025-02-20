package ffprobe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_process(t *testing.T) {
	t.Parallel()

	output := `{"format": {"bit_rate": "50"}, "streams": [{"index": 1}]}`

	actual, err := process(output)
	require.NoError(t, err)

	expected := &Result{
		Format: Format{
			Bitrate: "50",
		},
		Streams: []Stream{
			{
				Index: 1,
			},
		},
	}

	require.Equal(t, expected, actual)
}

func Test_process_ErrNoProbeResult(t *testing.T) {
	t.Parallel()

	actual, err := process("")
	require.Error(t, err)

	require.EqualError(t, err, "ffprobe returned no information")
	assert.Nil(t, actual)
}

func Test_process_ErrUnmarshalProbeResult(t *testing.T) {
	t.Parallel()

	actual, err := process("{")
	require.Error(t, err)

	require.EqualError(t, err, "unable to unmarshal probe result: unexpected end of JSON input")
	assert.Nil(t, actual)
}

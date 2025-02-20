package ffprobe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/media/ffprobe"
)

func TestProbe(t *testing.T) {
	t.Parallel()

	actual, err := ffprobe.Probe("./fixtures/h240.mp4")
	require.NoError(t, err)

	expected := &ffprobe.Result{
		Format: ffprobe.Format{
			Bitrate:    "1215656",
			Duration:   "2.048000",
			Filename:   "./fixtures/h240.mp4",
			FormatName: "mov,mp4,m4a,3gp,3g2,mj2",
			Size:       "311208",
			StartTime:  "0.000000",
		},
		Streams: []ffprobe.Stream{
			{
				Bitrate:          "1011384",
				CodecName:        "h264",
				CodecType:        "video",
				FrameRateAverage: "25/1",
				FrameRateLow:     "25/1",
				Height:           240,
				Index:            0,
				Profile:          "High",
				SampleRate:       "",
				Width:            426,
			},
			{
				Bitrate:          "194960",
				CodecName:        "aac",
				CodecType:        "audio",
				FrameRateAverage: "0/0",
				FrameRateLow:     "0/0",
				Height:           0,
				Index:            1,
				Profile:          "LC",
				SampleRate:       "48000",
				Width:            0,
			},
			{
				Bitrate:          "15",
				CodecName:        "",
				CodecType:        "data",
				FrameRateAverage: "12800/512",
				FrameRateLow:     "0/0",
				Height:           0,
				Index:            2,
				Profile:          "",
				SampleRate:       "",
				Width:            0,
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestProbe_ErrMissingFilename(t *testing.T) {
	t.Parallel()

	actual, err := ffprobe.Probe("")
	require.Error(t, err)

	require.EqualError(t, err, "filename is required")
	assert.Nil(t, actual)
}

func TestProbe_ErrRunningFFprobe(t *testing.T) {
	t.Parallel()

	actual, err := ffprobe.Probe("./fixtures/notfound.mov")
	require.Error(t, err)

	require.EqualError(t, err, "error returned when running ffprobe: ./fixtures/notfound.mov: No such file or directory\n")
	assert.Nil(t, actual)
}

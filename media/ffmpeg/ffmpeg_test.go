package ffmpeg_test

import (
	"testing"

	"github.com/go-cmd/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/io/filesystem"
	"github.com/sjdaws/pkg/media/ffmpeg"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	output := t.TempDir() + "/" + t.Name() + ".mp4"
	storage := filesystem.Default()

	command, err := ffmpeg.
		Input("./fixtures/h240.mp4").
		Output(output, ffmpeg.Option{Key: "vf", Value: "thumbnail"}).
		Command()
	require.NoError(t, err)

	assert.IsType(t, &cmd.Cmd{}, command)

	// Block until complete
	status := <-command.Start()

	require.NoError(t, status.Error)

	assert.True(t, storage.FileExists(output))
	assert.True(t, status.Complete)
	assert.Contains(t, status.Stdout, "frame=1")
	assert.Contains(t, status.Stdout, "out_time_ms=1240000")
	assert.Contains(t, status.Stdout, "progress=end")
}

func TestCommand_ErrGetCommandArguments(t *testing.T) {
	t.Parallel()

	command, err := ffmpeg.Input("").Output("").Command()
	require.Error(t, err)

	require.EqualError(t, err, "unable to determine command arguments: input filename is required")
	assert.Nil(t, command)
}

func TestRun(t *testing.T) {
	t.Parallel()

	output := t.TempDir() + "/" + t.Name() + ".mp4"
	storage := filesystem.Default()

	err := ffmpeg.
		Input("./fixtures/h240.mp4").
		Output(output, ffmpeg.Option{Key: "vf", Value: "thumbnail"}).
		Run()
	require.NoError(t, err)

	assert.True(t, storage.FileExists(output))
}

func TestRun_ErrGetCommandArguments(t *testing.T) {
	t.Parallel()

	err := ffmpeg.Input("").Output("").Run()
	require.Error(t, err)

	require.EqualError(t, err, "unable to determine command arguments: input filename is required")
}

func TestRun_ErrRunningFFmpeg(t *testing.T) {
	t.Parallel()

	output := t.TempDir() + "/" + t.Name() + ".mp4"

	err := ffmpeg.Input("./fixtures/notfound.mov").Output(output).Run()
	require.Error(t, err)

	require.ErrorContains(t, err, "error returned when running ffmpeg: ")
	require.ErrorContains(t, err, "Error opening input files: No such file or directory\n")
}

func TestSupports(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		codec    string
		expected bool
	}{
		"known supported": {
			codec:    "libx264",
			expected: true,
		},
		"known unsupported": {
			codec:    "this codec is not supported",
			expected: false,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := ffmpeg.Supports(testcase.codec)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}

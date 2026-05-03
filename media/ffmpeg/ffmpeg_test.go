package ffmpeg_test

import (
	"os/exec"
	"testing"

	"github.com/go-cmd/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sjdaws.com/pkg/io/filesystem"
	"sjdaws.com/pkg/media/ffmpeg"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	// Skip test if ffmpeg not installed
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("Skipping test: ffmpeg not installed")
	}

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
	assert.Contains(t, status.Stdout, "total_size=60952")
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

	// Skip test if ffmpeg not installed
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("Skipping test: ffmpeg not installed")
	}

	output := t.TempDir() + "/" + t.Name() + ".mp4"
	storage := filesystem.Default()

	err = ffmpeg.
		Input("./fixtures/h240.mp4").
		Output(output, ffmpeg.Option{Key: "vf", Value: "thumbnail"}).
		Run(t.Context())
	require.NoError(t, err)

	assert.True(t, storage.FileExists(output))
}

func TestRun_ErrGetCommandArguments(t *testing.T) {
	t.Parallel()

	err := ffmpeg.Input("").Output("").Run(t.Context())
	require.Error(t, err)

	require.EqualError(t, err, "unable to determine command arguments: input filename is required")
}

func TestRun_ErrRunningFFmpeg(t *testing.T) {
	t.Parallel()

	// Skip test if ffmpeg not installed
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("Skipping test: ffmpeg not installed")
	}

	output := t.TempDir() + "/" + t.Name() + ".mp4"

	err = ffmpeg.Input("./fixtures/notfound.mov").Output(output).Run(t.Context())
	require.Error(t, err)

	require.ErrorContains(t, err, "error returned when running ffmpeg: ")
	require.ErrorContains(t, err, "Error opening input files: No such file or directory\n")
}

func TestSupports(t *testing.T) {
	t.Parallel()

	// Skip test if ffmpeg not installed
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("Skipping test: ffmpeg not installed")
	}

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

			actual := ffmpeg.Supports(t.Context(), testcase.codec)

			assert.Equal(t, testcase.expected, actual)
		})
	}
}

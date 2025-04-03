package ffprobe

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/sjdaws/pkg/errors"
)

// Format result from ffprobe.
type Format struct {
	Bitrate    string `json:"bit_rate"`
	Duration   string `json:"duration"`
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	Size       string `json:"size"`
	StartTime  string `json:"start_time"`
}

// Result of a probe request.
type Result struct {
	Format  Format   `json:"format"`
	Streams []Stream `json:"streams"`
}

// Stream result from ffprobe.
type Stream struct {
	Bitrate          string `json:"bit_rate"`
	CodecName        string `json:"codec_name"`
	CodecType        string `json:"codec_type"`
	FrameRateAverage string `json:"avg_frame_rate"`
	FrameRateLow     string `json:"r_frame_rate"`
	Height           int    `json:"height"`
	Index            int    `json:"index"`
	Profile          string `json:"profile"`
	SampleRate       string `json:"sample_rate"`
	Width            int    `json:"width"`
}

// Probe a file and return metadata.
func Probe(filename string) (*Result, error) {
	if filename == "" {
		return nil, errors.New("filename is required")
	}

	args := []string{"-hide_banner", "-print_format", "json", "-show_format", "-show_streams", filename}

	var stdOut, stdErr strings.Builder

	//nolint:gosec // Assignment to variable intentional to overload stderr and stdout
	command := exec.Command("ffprobe", args...)
	command.Stderr = &stdErr
	command.Stdout = &stdOut

	err := command.Run()
	if err != nil {
		// This wrapping is intentionally backwards as we want the error rather than the exit code
		// ffprobe will return a single line as the error 'exit status XXX' which doesn't help.
		err = errors.Wrap(errors.New(stdErr.String()), err)

		return nil, errors.Wrap(err, "error returned when running ffprobe")
	}

	return process(stdOut.String())
}

func process(output string) (*Result, error) {
	if strings.TrimSpace(output) == "" {
		return nil, errors.New("ffprobe returned no information")
	}

	var result Result

	err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal probe result")
	}

	return &result, nil
}

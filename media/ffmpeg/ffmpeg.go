package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-cmd/cmd"

	"github.com/sjdaws/pkg/errors"
)

// Conversion ready to be passed to ffmpeg.
type Conversion struct {
	force  bool
	input  file
	output file
}

// Option for input or output.
type Option struct {
	Key   string
	Value any
}

// Request made to ffmpeg.
type Request struct {
	file
}

// file representation of an input or output.
type file struct {
	filename string
	options  []Option
}

// Input set input file and options.
func Input(filename string, options ...Option) *Request {
	return &Request{
		file: file{
			filename: filename,
			options:  options,
		},
	}
}

// Supports determine if a specific codec is supported.
func Supports(codec string) bool {
	var stdOut strings.Builder

	command := exec.Command("ffmpeg", "-codecs")
	command.Stdout = &stdOut

	// Errors don't matter here, if ffmpeg is missing there are bigger issues
	_ = command.Run()

	return strings.Contains(strings.ToLower(stdOut.String()), strings.ToLower(codec))
}

// Output set output file and options.
func (r *Request) Output(filename string, options ...Option) *Conversion {
	return &Conversion{
		force: false,
		input: file{
			filename: r.filename,
			options:  r.options,
		},
		output: file{
			filename: filename,
			options:  options,
		},
	}
}

// Command return a command which can run a conversion asynchronously.
func (c *Conversion) Command() (*cmd.Cmd, error) {
	arguments, err := c.getCommandArguments()
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine command arguments")
	}

	options := cmd.Options{
		BeforeExec:     nil,
		Buffered:       true,
		CombinedOutput: false,
		LineBufferSize: 0,
		Streaming:      false,
	}

	return cmd.NewCmdOptions(options, "ffmpeg", arguments...), nil
}

// Force a transformation to run, overwriting any existing files.
func (c *Conversion) Force() *Conversion {
	c.force = true

	return c
}

// Run a transformation.
func (c *Conversion) Run() error {
	arguments, err := c.getCommandArguments()
	if err != nil {
		return errors.Wrap(err, "unable to determine command arguments")
	}

	var stdErr strings.Builder

	command := exec.Command("ffmpeg", arguments...)
	command.Stderr = &stdErr

	err = command.Run()
	if err != nil {
		// This wrapping is intentionally backwards as we want the error rather than the exit code
		// ffmpeg will return a single line as the error 'exit status XXX' which doesn't help.
		err = errors.Wrap(errors.New(stdErr.String()), err)

		return errors.Wrap(err, "error returned when running ffmpeg")
	}

	return nil
}

// getCommandArguments for running ffmpeg.
func (c *Conversion) getCommandArguments() ([]string, error) {
	if strings.TrimSpace(c.input.filename) == "" {
		return nil, errors.New("input filename is required")
	}

	if strings.TrimSpace(c.output.filename) == "" {
		return nil, errors.New("output filename is required")
	}

	args := []string{"-hide_banner", "-v", "error"}
	args = append(args, c.processOptions(c.input.options)...)
	args = append(args, "-i", c.input.filename)
	args = append(args, c.processOptions(c.output.options)...)
	args = append(args, c.output.filename)
	args = append(args, "-progress", "/dev/stdout")

	if c.force {
		args = append(args, "-y")
	}

	return args, nil
}

// processOptions for input and output into key/value pairs.
func (c *Conversion) processOptions(options []Option) []string {
	result := make([]string, 0)

	for _, option := range options {
		if strings.TrimSpace(option.Key) == "" {
			continue
		}

		result = append(result, "-"+option.Key)

		if option.Value != nil {
			// Remove quotes from args
			if value, ok := option.Value.(string); ok {
				option.Value = strings.TrimSpace(strings.Trim(strings.TrimSpace(value), `"`))
			}

			if option.Value != "" {
				result = append(result, fmt.Sprintf("%v", option.Value))
			}
		}
	}

	return result
}

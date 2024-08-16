package initialize

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"

	"github.com/kilianpaquier/cli-sdk/pkg/clog"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// ErrAlreadyInitialized is the error returned (wrapped) when Run function is called but the project is already initialized.
var ErrAlreadyInitialized = errors.New("project already initialized")

// RunOption represents an option to be given to Run function.
type RunOption func(option) option

// WithLogger defines the logger implementation for Run function.
//
// When not provided, the default one used is the one from std log library.
func WithLogger(log clog.Logger) RunOption {
	return func(o option) option {
		o.log = log
		return o
	}
}

// WithReader defines the reader from where read answers ask during Run function.
//
// By default, when not provided, the reader used is os.Stdin (terminal).
//
// It can be useful in case of tests (see run_test.go for an example), but it's not recommended to use it in production code unless you're sure of what you're doing.
func WithReader(r io.Reader) RunOption {
	return func(o option) option {
		o.reader = r
		return o
	}
}

// Ask is the signature function for asking questions to the end user.
type Ask func(question string, args ...any) *string

// Warn is the signature function to print something in WARN level to the end user.
type Warn func(string, ...any)

// InputReader is the signature function for functions reading user inputs.
// Inspiration can be found with ReadMaintainer and ReadChart functions.
type InputReader func(log clog.Logger, config craft.Configuration, ask Ask) craft.Configuration

// WithInputReaders sets (it overrides the previously defined functions everytime it's called) the functions reading user inputs in Run function.
func WithInputReaders(inputs ...InputReader) RunOption {
	return func(o option) option {
		o.inputReaders = inputs
		return o
	}
}

// option represents the struct with all available options in Run function.
type option struct {
	inputReaders []InputReader
	log          clog.Logger
	reader       io.Reader
}

// newOpt creates a new option struct with all input Option functions while taking care of default values.
func newOpt(opts ...RunOption) option {
	o := option{}
	for _, opt := range opts {
		if opt != nil {
			o = opt(o)
		}
	}

	if len(o.inputReaders) == 0 {
		o.inputReaders = []InputReader{ReadMaintainer, ReadChart} // order is important since they will be executed in the same order
	}
	if o.log == nil {
		o.log = clog.Std()
	}
	if o.reader == nil {
		o.reader = os.Stdin
	}

	return o
}

// Run initializes a new craft project in case a craft.CraftFile doesn't exist in destdir.
// All user inputs must be configured through WithInputReaders option, by default the main maintainer and chart generation will be asked.
//
// Multiple options can be given like saving craft configuration file at the end (default is false),
// the logger used to ask question to the end user
// or the reader from where retrieve the end user answers (but as provided in WithReader doc it should be used with caution - you should know what you're doing).
//
// In case craft.CraftFile already exists, it's read and returned alongside ErrAlreadyInitialized (should be handled in caller).
func Run(_ context.Context, destdir string, opts ...RunOption) (craft.Configuration, error) {
	o := newOpt(opts...)

	// read config configuration
	var config craft.Configuration
	err := craft.Read(destdir, &config)
	if err == nil {
		return config, ErrAlreadyInitialized
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return craft.Configuration{}, fmt.Errorf("%s exists but is not readable: %w", craft.File, err)
	}

	// keep scanner and log abstracted for better developer friendly
	// and also avoid potential issues (scanner reaffected nil in input function)
	scanner := bufio.NewScanner(o.reader)
	ask := func(question string, args ...any) *string {
		o.log.Infof(question, args...)
		scanner.Scan()
		if answer := scanner.Text(); answer != "" {
			return &answer
		}
		return nil
	}
	for _, inputs := range o.inputReaders {
		config = inputs(o.log, config, ask) // read all configured inputs (default is main maintainer and chart generation)
	}

	return config, nil
}

// ReadMaintainer creates a maintainer with Q&A method from the end user.
func ReadMaintainer(_ clog.Logger, config craft.Configuration, ask Ask) craft.Configuration {
	var maintainer craft.Maintainer

	// main maintainer name
	var name *string
	for name == nil {
		name = ask("Who's the maintainer name (required, it can be a group name, anything) ?")
	}
	maintainer.Name = *name
	maintainer.Email = ask("Who's the maintainer email (optional, press Enter to skip) ?")
	maintainer.URL = ask("Who's the maintainer url (optional, press Enter to skip) ?")

	config.Maintainers = append(config.Maintainers, maintainer)
	return config
}

// ReadChart retrieves the chart generation choice from the end user.
func ReadChart(log clog.Logger, config craft.Configuration, ask Ask) craft.Configuration {
	// Helm chart generation
	for {
		chart := ask("Would you like to generate an Helm chart (optional, press Enter to skip, default is truthy) 0/1 ?")
		if chart == nil {
			return config // no response provided, going through with chart activated
		}

		value, err := strconv.ParseBool(*chart)
		if err != nil {
			log.Warnf("invalid chart answer '%s', must be a boolean", *chart)
			continue
		}

		config.NoChart = !value
		return config
	}
}

package generate

import (
	"os"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
)

// ExecOpts represents all options given to ExecFunc functions.
type ExecOpts struct {
	FileHandlers []FileHandler

	Force    []string
	ForceAll bool

	EndDelim   string
	StartDelim string
}

// RunOption is the right function to tune Run function with specific behaviors.
type RunOption func(runOptions) runOptions

// WithMetaHandlers is an option for Run function.
// It specifies the slice of MetaHandler, which defines the behavior for files and directories generation.
//
// When not given, MetaHandlers' function result is used as default slice.
func WithMetaHandlers(handlers ...MetaHandler) RunOption {
	return func(o runOptions) runOptions {
		o.metaHandlers = handlers
		return o
	}
}

// WithDelimiters is an option for Run function to use specific go template delimiters.
//
// If not given, default delimiters are << and >>.
func WithDelimiters(startDelim, endDelim string) RunOption {
	return func(o runOptions) runOptions {
		o.startDelim = startDelim
		o.endDelim = endDelim
		return o
	}
}

// WithDestination is an option for Run function to specify
// the destination directory of generation.
//
// If not given, default destination is the current directory where Run is executed.
func WithDestination(destdir string) RunOption {
	return func(o runOptions) runOptions {
		o.destdir = &destdir
		return o
	}
}

// WithDetects is an option for Run function defining the detections (languages) to identify.
//
// When not given, Detects is used as default slice.
func WithDetects(funcs ...DetectFunc) RunOption {
	return func(o runOptions) runOptions {
		o.detectFuncs = funcs
		return o
	}
}

// WithForce is an option for Run function to specify which
// files must be generated even if the top notice is not present anymore (see IsGenerated).
//
// If not given, no files are force'd generated.
func WithForce(filenames ...string) RunOption {
	return func(o runOptions) runOptions {
		o.force = filenames
		return o
	}
}

// WithForceAll is an option for Run function to specify
// whether to force the generation of all files or not.
// When given, WithForce isn't used.
//
// If not given, this option is false.
func WithForceAll(forceAll bool) RunOption {
	return func(o runOptions) runOptions {
		o.forceAll = forceAll
		return o
	}
}

// WithTemplates is an option for Run function to specify the templates directory and filesystem.
//
// Please not that the input dir path separator must be the one used with path.Join
// and not the one OS specific from filepath.Join.
//
// If not given, default filesystem is the embedded one FS.
func WithTemplates(dir string, fs cfs.FS) RunOption {
	return func(o runOptions) runOptions {
		o.tmplDir = dir
		o.fs = fs
		return o
	}
}

// runOptions is the struct related to Option function(s) defining all optional properties.
type runOptions struct {
	detectFuncs  []DetectFunc
	metaHandlers []MetaHandler

	destdir *string

	force    []string
	forceAll bool

	fs      cfs.FS
	tmplDir string

	endDelim   string
	startDelim string
}

// newRunOpt creates a new option struct with all input Option functions
// while taking care of default values.
func newRunOpt(opts ...RunOption) runOptions {
	var ro runOptions
	for _, opt := range opts {
		if opt != nil {
			ro = opt(ro)
		}
	}

	if ro.startDelim == "" || ro.endDelim == "" {
		ro.startDelim = "<<"
		ro.endDelim = ">>"
	}
	if ro.destdir == nil {
		dir, _ := os.Getwd()
		ro.destdir = &dir
	}
	if ro.fs == nil {
		ro.fs = FS()
		ro.tmplDir = "templates"
	}
	if len(ro.detectFuncs) == 0 {
		ro.detectFuncs = DetectFuncs()
	}
	if len(ro.metaHandlers) == 0 {
		ro.metaHandlers = MetaHandlers()
	}

	return ro
}

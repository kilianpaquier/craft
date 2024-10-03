package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"golang.org/x/mod/modfile"

	"github.com/kilianpaquier/craft/pkg/craft"
)

var (
	errMissingModuleStatement = errors.New("invalid go.mod, module statement is missing")
	errMissingGoStatement     = errors.New("invalid go.mod, go statement is missing")
)

var versionRegexp = regexp.MustCompile("^v[0-9]+$")

// Gomod represents the parsed struct for go.mod file
type Gomod struct {
	LangVersion string
	Platform    string
	ProjectHost string
	ProjectName string
	ProjectPath string
}

// DetectGolang handles the detection of golang at destdir.
//
// A valid golang project must have a valid go.mod file.
func DetectGolang(ctx context.Context, log clog.Logger, destdir string, metadata *Metadata) ([]Exec, error) {
	gomod := filepath.Join(destdir, craft.Gomod)
	gocmd := filepath.Join(destdir, craft.Gocmd)

	// retrieve module from go.mod
	statements, err := readGomod(gomod)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("read go.mod: %w", err)
		}
		return nil, nil
	}

	metadata.Platform = statements.Platform
	metadata.ProjectHost = statements.ProjectHost
	metadata.ProjectName = statements.ProjectName
	metadata.ProjectPath = statements.ProjectPath

	// check hugo detection
	if execs, _ := detectHugo(ctx, log, destdir, metadata); len(execs) > 0 { //nolint:revive
		return execs, nil
	}

	log.Infof("golang detected, a %s is present and valid", craft.Gomod)
	metadata.Languages["golang"] = statements

	entries, err := os.ReadDir(gocmd)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Warnf("failed to read directory: %s", err.Error())
	}

	// range over folders to retrieve binaries type
	for _, entry := range entries {
		if entry.IsDir() {
			switch {
			case strings.HasPrefix(entry.Name(), "cron-"):
				metadata.Crons[entry.Name()] = struct{}{}
			case strings.HasPrefix(entry.Name(), "job-"):
				metadata.Jobs[entry.Name()] = struct{}{}
			case strings.HasPrefix(entry.Name(), "worker-"):
				metadata.Workers[entry.Name()] = struct{}{}
			default:
				// by default, executables in cmd folder are CLI
				metadata.Clis[entry.Name()] = struct{}{}
			}
			metadata.Binaries++
		}
	}

	return []Exec{DefaultExec("lang_golang")}, nil
}

var _ Detect = DetectGolang // ensure interface is implemented

// detectHugo handles the detection of hugo at destdir.
func detectHugo(_ context.Context, log clog.Logger, destdir string, metadata *Metadata) ([]Exec, error) {
	// detect hugo project
	configs, _ := filepath.Glob(filepath.Join(destdir, "hugo.*"))

	// detect hugo theme
	themes, _ := filepath.Glob(filepath.Join(destdir, "theme.*"))

	if len(configs) > 0 || len(themes) > 0 {
		log.Infof("hugo detected, a hugo configuration file or hugo theme file is present")

		if metadata.CI != nil {
			log.Warnf("removing codecov, codeql and sonar option from CI since they're not available with hugo projects")
			metadata.CI.Options = slices.DeleteFunc(metadata.CI.Options, func(option string) bool {
				return slices.Contains([]string{craft.CodeCov, craft.CodeQL, craft.Sonar}, option)
			})
		}

		metadata.Languages["hugo"] = nil
		return []Exec{DefaultExec("lang_hugo")}, nil
	}
	return nil, nil
}

var _ Detect = detectHugo // ensure interface is implemented

// readGomod reads the go.mod file at modpath input and returns its gomod representation.
func readGomod(modpath string) (Gomod, error) {
	// read go.mod at modpath
	bytes, err := os.ReadFile(modpath)
	if err != nil {
		return Gomod{}, fmt.Errorf("read file: %w", err)
	}

	// parse go.mod into it's modfile representation
	file, err := modfile.Parse(modpath, bytes, nil)
	if err != nil {
		return Gomod{}, fmt.Errorf("parse go.mod: %w", err)
	}

	var errs []error
	var gomod Gomod

	// parse module statement
	if file.Module == nil || file.Module.Mod.Path == "" {
		errs = append(errs, errMissingModuleStatement)
	} else {
		gomod.ProjectHost, gomod.ProjectPath = func() (host, subpath string) {
			sections := strings.Split(file.Module.Mod.Path, "/")
			if versionRegexp.MatchString(sections[len(sections)-1]) {
				return sections[0], strings.Join(sections[1:len(sections)-1], "/") // retrieve all sections but the last element
			}
			return sections[0], strings.Join(sections[1:], "/") // retrieve all sections
		}()
		gomod.Platform, _ = ParsePlatform(gomod.ProjectHost)
		gomod.ProjectName = path.Base(gomod.ProjectPath)
	}

	// parse go statement
	if file.Go == nil {
		errs = append(errs, errMissingGoStatement)
	} else {
		gomod.LangVersion = file.Go.Version
	}

	// override lang version if toolchain is present
	// it's preempting provided go version for build purposes
	if file.Toolchain != nil {
		gomod.LangVersion = file.Toolchain.Name[2:]
	}

	return gomod, errors.Join(errs...)
}

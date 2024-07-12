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

	"golang.org/x/mod/modfile"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/logger"
)

var _versionRegexp = regexp.MustCompile("^v[0-9]+$")

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
func DetectGolang(ctx context.Context, log logger.Logger, destdir string, metadata Metadata) (Metadata, []Exec) {
	gomod := filepath.Join(destdir, craft.Gomod)
	gocmd := filepath.Join(destdir, craft.Gocmd)

	// retrieve module from go.mod
	statements, err := readGomod(gomod)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Warnf("failed to parse %s statements: %s", craft.Gomod, err.Error())
		}
		return metadata, nil
	}

	metadata.Platform = statements.Platform
	metadata.ProjectHost = statements.ProjectHost
	metadata.ProjectName = statements.ProjectName
	metadata.ProjectPath = statements.ProjectPath

	// check hugo detection
	if metadata, exec := detectHugo(ctx, log, destdir, metadata); len(exec) > 0 { // nolint:revive
		return metadata, exec
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

	return metadata, []Exec{DefaultExec("lang_golang")}
}

var _ Detect = DetectGolang // ensure interface is implemented

// detectHugo handles the detection of hugo at destdir.
func detectHugo(_ context.Context, log logger.Logger, destdir string, metadata Metadata) (Metadata, []Exec) {
	// detect hugo project
	configs, _ := filepath.Glob(filepath.Join(destdir, "hugo.*"))

	// detect hugo theme
	themes, _ := filepath.Glob(filepath.Join(destdir, "theme.*"))

	if len(configs) > 0 || len(themes) > 0 {
		log.Info("hugo detected, a hugo configuration file or hugo theme file is present")

		if metadata.CI != nil {
			if slices.Contains(metadata.CI.Options, craft.CodeQL) {
				log.Warn("codeql option is not available with hugo generation, deactivating it")
				metadata.CI.Options = slices.DeleteFunc(metadata.CI.Options, func(option string) bool {
					return option == craft.CodeQL
				})
			}

			if slices.Contains(metadata.CI.Options, craft.CodeCov) {
				log.Warn("codecov option is not available with hugo generation, deactivating it")
				metadata.CI.Options = slices.DeleteFunc(metadata.CI.Options, func(option string) bool {
					return option == craft.CodeCov
				})
			}
		}

		metadata.Languages["hugo"] = nil
		return metadata, []Exec{DefaultExec("lang_hugo")}
	}
	return metadata, nil
}

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
		errs = append(errs, errors.New("invalid go.mod, module statement is missing"))
	} else {
		gomod.ProjectHost, gomod.ProjectPath = func() (host, subpath string) {
			sections := strings.Split(file.Module.Mod.Path, "/")
			if _versionRegexp.MatchString(sections[len(sections)-1]) {
				return sections[0], strings.Join(sections[1:len(sections)-1], "/") // retrieve all sections but the last element
			}
			return sections[0], strings.Join(sections[1:], "/") // retrieve all sections
		}()
		gomod.Platform, _ = ParsePlatform(gomod.ProjectHost)
		gomod.ProjectName = path.Base(gomod.ProjectPath)
	}

	// parse go statement
	if file.Go == nil {
		errs = append(errs, errors.New("invalid go.mod, go statement is missing"))
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

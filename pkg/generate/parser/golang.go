package parser

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

const (
	// FolderCMD represents the cmd folder where go main.go should be placed according to go layout.
	FolderCMD = "cmd"
	// FileGomod represents the go.mod filename.
	FileGomod = "go.mod"
)

var (
	// ErrMissingModuleStatement is the error returned when module statement is missing from go.mod.
	ErrMissingModuleStatement = errors.New("invalid go.mod, module statement is missing")

	// ErrMissingGoStatement is the error returned when go statement is missing from go.mod.
	ErrMissingGoStatement = errors.New("invalid go.mod, go statement is missing")
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

// Golang handles the parsing of a golang repository at destdir.
//
// A valid golang project must have a valid go.mod file.
func Golang(ctx context.Context, destdir string, config *craft.Config) error {
	gomod := filepath.Join(destdir, FileGomod)
	gocmd := filepath.Join(destdir, FolderCMD)

	// retrieve module from go.mod
	statements, err := readGomod(gomod)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("read %s: %w", FileGomod, err)
		}
		return nil
	}

	config.Platform = statements.Platform
	config.ProjectHost = statements.ProjectHost
	config.ProjectName = statements.ProjectName
	config.ProjectPath = statements.ProjectPath

	// check hugo repository
	if ok := isHugo(ctx, destdir, config); ok {
		return nil
	}

	generate.GetLogger(ctx).Infof("golang detected, file '%s' is present and valid", FileGomod)
	config.SetLanguage("golang", statements)

	entries, err := os.ReadDir(gocmd)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		generate.GetLogger(ctx).Warnf("failed to read directory: %s", err.Error())
	}

	// range over folders to retrieve binaries type
	for _, entry := range entries {
		if entry.IsDir() {
			switch {
			case strings.HasPrefix(entry.Name(), "cron-"):
				config.SetCron(entry.Name())
			case strings.HasPrefix(entry.Name(), "job-"):
				config.SetJob(entry.Name())
			case strings.HasPrefix(entry.Name(), "worker-"):
				config.SetWorker(entry.Name())
			default:
				// by default, executables in cmd folder are CLI
				config.SetCLI(entry.Name())
			}
		}
	}
	return nil
}

var _ generate.Parser[craft.Config] = Golang // ensure interface is implemented

func isHugo(ctx context.Context, destdir string, config *craft.Config) bool {
	// detect hugo project
	configs, _ := filepath.Glob(filepath.Join(destdir, "hugo.*"))

	// detect hugo theme
	themes, _ := filepath.Glob(filepath.Join(destdir, "theme.*"))

	if len(configs) > 0 || len(themes) > 0 {
		config.SetLanguage("hugo", nil)
		generate.GetLogger(ctx).Infof("hugo detected, a hugo configuration file or hugo theme file is present")
		return true
	}
	return false
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
		errs = append(errs, ErrMissingModuleStatement)
	} else {
		gomod.ProjectHost, gomod.ProjectPath = func() (host, subpath string) {
			sections := strings.Split(file.Module.Mod.Path, "/")
			if versionRegexp.MatchString(sections[len(sections)-1]) {
				return sections[0], strings.Join(sections[1:len(sections)-1], "/") // retrieve all sections but the last element
			}
			return sections[0], strings.Join(sections[1:], "/") // retrieve all sections
		}()
		gomod.Platform, _ = parsePlatform(gomod.ProjectHost)
		gomod.ProjectName = path.Base(gomod.ProjectPath)
	}

	// parse go statement
	if file.Go == nil {
		errs = append(errs, ErrMissingGoStatement)
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

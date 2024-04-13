package detectgen

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

	"github.com/sirupsen/logrus"
	"golang.org/x/mod/modfile"

	"github.com/kilianpaquier/craft/internal/generate/remote"
	"github.com/kilianpaquier/craft/internal/models"
)

// detectGolang handles the detection of golang at config provided destination directory.
//
// It returns the appropriate slice of GenerateFunc depending on golang's craft related options (hugo).
// A valid golang project must have a valid go.mod file.
//
// The detection also handles parsing executables present in cmd folder.
func detectGolang(ctx context.Context, config *models.GenerateConfig) []GenerateFunc {
	log := logrus.WithContext(ctx)

	gomod := filepath.Join(config.Options.DestinationDir, models.Gomod)
	gocmd := filepath.Join(config.Options.DestinationDir, models.Gocmd)

	// retrieve module from go.mod
	statements, err := readGomod(gomod)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.WithError(err).Warnf("failed to parse %s statements", models.Gomod)
		}
		return nil
	}

	config.LangVersion = statements.LangVersion
	config.Platform = statements.Platform
	config.ProjectHost = statements.ProjectHost
	config.ProjectName = statements.ProjectName
	config.ProjectPath = statements.ProjectPath

	// check hugo detection
	if generates := detectHugo(ctx, config); len(generates) > 0 {
		return generates
	}

	log.Infof("golang detected, a %s is present and valid", models.Gomod)
	config.Languages = append(config.Languages, string(NameGolang))

	entries, err := os.ReadDir(gocmd)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Debugf("%s doesn't exist", gocmd)
		} else {
			log.WithError(err).Warnf("failed to read %s folder", gocmd)
		}
	}

	// range over folders to retrieve binaries type
	for _, entry := range entries {
		if entry.IsDir() {
			switch {
			case strings.HasPrefix(entry.Name(), "cron-"):
				config.Crons[entry.Name()] = struct{}{}
			case strings.HasPrefix(entry.Name(), "job-"):
				config.Jobs[entry.Name()] = struct{}{}
			case strings.HasPrefix(entry.Name(), "worker-"):
				config.Workers[entry.Name()] = struct{}{}
			default:
				// by default, executables in cmd folder are CLI
				config.Clis[entry.Name()] = struct{}{}
			}
			config.Binaries++
		}
	}

	return []GenerateFunc{GetGenerateFunc(NameGolang)}
}

// gomod represents the parsed struct for go.mod file
type gomod struct {
	LangVersion string
	Platform    string
	ProjectHost string
	ProjectName string
	ProjectPath string
}

// readGomod reads the go.mod file at modpath input and returns its gomod representation.
func readGomod(modpath string) (*gomod, error) {
	// read go.mod at modpath
	bytes, err := os.ReadFile(modpath)
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}

	// parse go.mod into it's modfile representation
	file, err := modfile.Parse(modpath, bytes, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}

	var errs []error
	gomod := &gomod{}

	// parse module statement
	if file.Module == nil || file.Module.Mod.Path == "" {
		errs = append(errs, errors.New("invalid go.mod, module statement is missing"))
	} else {
		gomod.ProjectHost, gomod.ProjectPath = func() (host, subpath string) {
			sections := strings.Split(file.Module.Mod.Path, "/")
			if regexp.MustCompile("^v[0-9]+$").MatchString(sections[len(sections)-1]) {
				return sections[0], strings.Join(sections[1:len(sections)-1], "/") // retrieve all sections but the last element
			}
			return sections[0], strings.Join(sections[1:], "/") // retrieve all sections
		}()
		gomod.Platform, _ = remote.ParsePlatform(gomod.ProjectHost)
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

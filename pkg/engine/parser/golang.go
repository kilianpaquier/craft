package parser

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

const (
	// FolderCMD represents the cmd folder where go main.go should be placed according to go layout.
	FolderCMD = "cmd"

	// FileGomod represents the go.mod filename.
	FileGomod = "go.mod"

	// FileMain represents the main.go filename.
	FileMain = "main.go"
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
	ModulePath  string
}

// AsVCS returns the vcs configuration associated to module statement in go.mod.
func (g Gomod) AsVCS() VCS {
	sections := strings.Split(g.ModulePath, "/")
	projectPath := func() string {
		if versionRegexp.MatchString(sections[len(sections)-1]) {
			return strings.Join(sections[1:len(sections)-1], "/") // retrieve all sections but the last element
		}
		return strings.Join(sections[1:], "/") // retrieve all sections
	}()

	return VCS{
		Platform:    func() string { p, _ := parsePlatform(sections[0]); return p }(),
		ProjectHost: sections[0],
		ProjectPath: projectPath,
		ProjectName: path.Base(projectPath),
	}
}

// ReadGomod reads the go.mod file at destdir
// and returns its representation.
//
// It will return an error if the go.mod file is missing the following properties:
//   - module statement
//   - go statement
//
// It can be used as a simple function, calling it directly,
// but can also be used as its expected usage with engine.Generate:
//
//	type config struct { ... }
//
//	func Golang(ctx context.Context, destdir string, c *config) error {
//		gomod, err := parser.ReadGomod(destdir)
//		if err != nil {
//			if errors.Is(err, fs.ErrNotExist) {
//				return nil
//			}
//			return fmt.Errorf("read go.mod: %w", err)
//		}
//		engine.GetLogger(ctx).Infof("golang detected, file '%s' is present and valid", parser.FileGomod)
//		// do something with gomod (e.g. update config since it's a pointer)
//		return nil
//	}
func ReadGomod(destdir string) (Gomod, error) {
	modpath := filepath.Join(destdir, FileGomod)

	// read go.mod
	bytes, err := os.ReadFile(modpath)
	if err != nil {
		return Gomod{}, fmt.Errorf("read file: %w", err)
	}

	// parse go.mod into it's modfile representation
	file, err := modfile.Parse(FileGomod, bytes, nil)
	if err != nil {
		return Gomod{}, fmt.Errorf("parse modfile: %w", err)
	}
	var gomod Gomod

	// parse module statement
	if file.Module == nil || file.Module.Mod.Path == "" {
		return Gomod{}, ErrMissingModuleStatement
	}
	gomod.ModulePath = file.Module.Mod.Path

	// parse go statement
	if file.Go == nil {
		return Gomod{}, ErrMissingGoStatement
	}
	gomod.LangVersion = file.Go.Version

	// override lang version if toolchain is present
	// it's preempting provided go version for build purposes
	if file.Toolchain != nil {
		gomod.LangVersion = file.Toolchain.Name[2:]
	}

	return gomod, nil
}

// ReadGoCmd reads cmd folder at destdir
// and returns all executables found.
//
// An executable is a main.go present in a subfolder of cmd.
// For instance cmd/<name>/main.go.
//
// Executables are categorized as CLI, Cron, Job and Worker.
// The prefix of <name> indicates in which category the executable belongs:
//   - <name> for CLI executables
//   - cron-<name> for cron executables
//   - job-<name> for job executables
//   - worker-<name> for worker executables
//
// It can be used as a simple function, calling it directly,
// but can also be used as its expected usage with engine.Generate:
//
//	type config struct { ... }
//
//	func Golang(ctx context.Context, destdir string, c *config) error {
//		gomod, err := parser.ReadGomod(destdir)
//		if err != nil {
//			if errors.Is(err, fs.ErrNotExist) {
//				return nil
//			}
//			return fmt.Errorf("read go.mod: %w", err)
//		}
//
//		engine.GetLogger(ctx).Infof("golang detected, file '%s' is present and valid", parser.FileGomod)
//		// do something with gomod (e.g. update config since it's a pointer)
//
//		executables, err := parser.ReadGoCmd(destdir)
//		if err != nil {
//			if errors.Is(err, fs.ErrNotExist) {
//				return nil
//			}
//			return fmt.Errorf("read cmd: %w", err)
//		}
//		// do something with executables (e.g. update config since it's a pointer)
//		return nil
//	}
func ReadGoCmd(destdir string) (Executables, error) {
	cmdpath := filepath.Join(destdir, FolderCMD)

	entries, err := os.ReadDir(cmdpath)
	if err != nil {
		return Executables{}, fmt.Errorf("read dir: %w", err)
	}

	// range over folders to retrieve binaries type
	var executables Executables
	for _, entry := range entries {
		if entry.IsDir() {
			if !files.Exists(filepath.Join(cmdpath, entry.Name(), FileMain)) {
				continue
			}

			switch {
			case strings.HasPrefix(entry.Name(), "cron-"):
				executables.SetCron(entry.Name())
			case strings.HasPrefix(entry.Name(), "job-"):
				executables.SetJob(entry.Name())
			case strings.HasPrefix(entry.Name(), "worker-"):
				executables.SetWorker(entry.Name())
			default:
				// by default, executables in cmd folder are CLI
				executables.SetCLI(entry.Name())
			}
		}
	}
	return executables, nil
}

// HugoConfig represents the parse'd hugo.* or theme.* file associated to hugo configuration file.
type HugoConfig struct{}

// Hugo detects if the project is a Hugo project.
//
// Detection consists of looking for hugo.* or theme.* files in the given destdir.
//
// It can be used as a simple function, calling it directly,
// but can also be used as its expected usage with engine.Generate:
//
//	type config struct { ... }
//
//	func Hugo(ctx context.Context, destdir string, c *config) error {
//		hugoc, ok := parser.Hugo(destdir)
//		if !ok {
//			return nil
//		}
//		engine.GetLogger(ctx).Infof("hugo detected, theme or hugo files are present")
//		// do something with hugo config (e.g. update config since it's a pointer)
//		return nil
//	}
func Hugo(destdir string) (HugoConfig, bool) {
	// detect hugo project
	configs, _ := filepath.Glob(filepath.Join(destdir, "hugo.*"))

	// detect hugo theme
	themes, _ := filepath.Glob(filepath.Join(destdir, "theme.*"))

	return HugoConfig{}, len(configs) > 0 || len(themes) > 0
}

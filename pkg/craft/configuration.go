package craft

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"gopkg.in/yaml.v3"
)

// Configuration represents all options configurable in .craft file at root project.
//
// yaml tags are for .craft file and json tags for templating.
type Configuration struct {
	CI           *CI          `json:"-"                     yaml:"ci,omitempty"                             validate:"omitempty,required"`
	Description  *string      `json:"description,omitempty" yaml:"description,omitempty"`
	Docker       *Docker      `json:"docker,omitempty"      yaml:"docker,omitempty"                         validate:"omitempty,required"`
	License      *string      `json:"-"                     yaml:"license,omitempty"                        validate:"omitempty,oneof=agpl-3.0 apache-2.0 bsd-2-clause bsd-3-clause bsl-1.0 cc0-1.0 epl-2.0 gpl-2.0 gpl-3.0 lgpl-2.1 mit mpl-2.0 unlicense"`
	Maintainers  []Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"   builder:"append" validate:"required,dive,required"`
	NoChart      bool         `json:"-"                     yaml:"no_chart,omitempty"`
	NoGoreleaser bool         `json:"-"                     yaml:"no_goreleaser,omitempty"`
	NoMakefile   bool         `json:"-"                     yaml:"no_makefile,omitempty"`
	Platform     string       `json:"-"                     yaml:"platform,omitempty"                       validate:"omitempty,oneof=bitbucket gitea github gitlab"`
}

// CI is the struct for craft continuous integration tuning.
type CI struct {
	Name    string   `json:"-" yaml:"name,omitempty"                     validate:"required,oneof=github gitlab"`
	Options []string `json:"-" yaml:"options,omitempty" builder:"append" validate:"omitempty,dive,oneof=codecov codeql dependabot netlify pages renovate sonar"`
	Release Release  `json:"-" yaml:"release,omitempty"                  validate:"required"`
}

// Release is the struct for craft continuous integration release specifics configuration.
type Release struct {
	Action    string `json:"-" yaml:"action,omitempty"  validate:"omitempty,oneof=gh-release release-drafter semantic-release"`
	Auto      bool   `json:"-" yaml:"auto"`
	Backmerge bool   `json:"-" yaml:"backmerge"`
	Disable   bool   `json:"-" yaml:"disable,omitempty"`
	Mode      string `json:"-" yaml:"mode,omitempty"    validate:"omitempty,oneof=github-apps github-token personal-token"`
}

// Docker is the struct for craft docker tuning.
type Docker struct {
	Registry *string `json:"registry,omitempty" yaml:"registry,omitempty"`
	Port     *uint16 `json:"port,omitempty"     yaml:"port,omitempty"`
}

// Maintainer represents a project maintainer. It's inspired from helm Maintainer struct.
//
// The only difference are the present tags and the pointers on both email and url properties.
type Maintainer struct {
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`
	URL   *string `json:"url,omitempty"   yaml:"url,omitempty"`
	Name  string  `json:"name,omitempty"  yaml:"name,omitempty"  validate:"required"`
}

// Read reads the .craft file in srcdir input into the out input.
func Read(srcdir string, out any) error {
	src := filepath.Join(srcdir, File)

	content, err := os.ReadFile(src)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fs.ErrNotExist
		}
		return fmt.Errorf("read file: %w", err)
	}

	if err := yaml.Unmarshal(content, out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

// Write writes the input craft into the input destdir in .craft file.
func Write(destdir string, config Configuration) error {
	dest := filepath.Join(destdir, File)

	// create a buffer with craft notice
	buffer := bytes.NewBuffer([]byte("# Craft configuration file (https://github.com/kilianpaquier/craft)\n---\n"))

	// create yaml encoder and writes the full configuration in the buffer,
	// following the craft notice
	encoder := yaml.NewEncoder(buffer)
	defer encoder.Close()
	encoder.SetIndent(2)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("encode file: %w", err)
	}

	if err := os.WriteFile(dest, buffer.Bytes(), cfs.RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

// EnsureDefaults acts to ensure default properties are always sets
// and migrates old properties into new fields.
func (c Configuration) EnsureDefaults() Configuration {
	if c.CI != nil {
		// generic function to match an option included in a slice of options
		del := func(options ...string) func(option string) bool {
			return func(option string) bool {
				return slices.Contains(options, option)
			}
		}

		// migrate old auto_release option
		if slices.Contains(c.CI.Options, "auto_release") {
			c.CI.Release.Auto = true
			c.CI.Options = slices.DeleteFunc(c.CI.Options, del("auto_release"))
		}

		// migrate old backmerge option
		if slices.Contains(c.CI.Options, "backmerge") {
			c.CI.Release.Backmerge = true
			c.CI.Options = slices.DeleteFunc(c.CI.Options, del("backmerge"))
		}

		// set default release action in case it's not provided
		if c.CI.Release.Action == "" {
			c.CI.Release.Action = SemanticRelease
		}

		if c.CI.Name == Github {
			// ses default release mode for github actions
			if c.CI.Release.Mode == "" {
				c.CI.Release.Mode = GithubToken
			}

			if c.CI.Release.Action == ReleaseDrafter || c.CI.Release.Action == GhRelease {
				// remove backmerge feature on gh-release or release-drafter releasing since isn't not available
				c.CI.Release.Backmerge = false
				// set release mode to github-token since it's the only mode available with gh-release or release-drafter
				c.CI.Release.Mode = GithubToken
			}
		}

		if c.CI.Name == Gitlab {
			// keep release mode empty when working with gitlab CICD
			c.CI.Release.Mode = ""
		}
	}
	return c
}

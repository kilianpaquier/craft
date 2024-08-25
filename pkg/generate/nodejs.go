package generate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"

	"github.com/kilianpaquier/craft/pkg/craft"
)

var _packageManagerRegexp = regexp.MustCompile(`^(npm|pnpm|yarn|bun)@\d+\.\d+\.\d+(-.+)?$`)

// PackageJSON represents the node package json descriptor.
type PackageJSON struct {
	Author         *string  `json:"author,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Files          []string `json:"files,omitempty"`
	Keywords       []string `json:"keywords,omitempty"`
	License        *string  `json:"license,omitempty"`
	Main           *string  `json:"main,omitempty"`
	Module         string   `json:"module,omitempty"`
	Name           string   `json:"name,omitempty"           validate:"required"`
	PackageManager string   `json:"packageManager,omitempty" validate:"required"`
	Private        bool     `json:"private,omitempty"`
	PublishConfig  struct {
		Access     string `json:"access,omitempty"`
		Provenance bool   `json:"provenance,omitempty"`
		Registry   string `json:"registry,omitempty"`
		Tag        string `json:"tag,omitempty"`
	} `json:"publishConfig,omitempty"`
	Repository *struct {
		URL string `json:"url,omitempty" validate:"required"`
	} `json:"repository,omitempty" validate:"required_if=Private false"`
	Scripts map[string]string `json:"scripts,omitempty"`
	Version string            `json:"version,omitempty"`
}

// Validate validates the given PackageJSON struct.
func (p *PackageJSON) Validate() error {
	var errs []error

	if p.PackageManager != "" && !_packageManagerRegexp.MatchString(p.PackageManager) {
		// json schema takes care of saying which regexp must be validated
		errs = append(errs, errors.New("package.json packageManager isn't valid"))
	}

	if err := validator.New().Struct(p); err != nil {
		errs = append(errs, fmt.Errorf("struct validation: %w", err))
	}
	return errors.Join(errs...)
}

// DetectNodejs handles nodejs detection at destdir.
// It scans the project for a package.json and validates it.
func DetectNodejs(_ context.Context, log clog.Logger, destdir string, metadata Metadata) (Metadata, []Exec, error) {
	jsonpath := filepath.Join(destdir, craft.PackageJSON)
	pkg, err := readPackageJSON(jsonpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return metadata, nil, fmt.Errorf("read package.json: %w", err)
		}
		return metadata, nil, nil
	}

	log.Infof("nodejs detected, a '%s' is present and valid", craft.PackageJSON)

	metadata.Languages["nodejs"] = pkg
	metadata.ProjectName = pkg.Name
	if pkg.Main != nil {
		metadata.Binaries++
	}

	// deactivate makefile because commands are facilitated by package.json scripts
	if !metadata.NoMakefile {
		log.Warnf("makefile option not available with nodejs generation, deactivating it")
		metadata.NoMakefile = true
	}

	return metadata, []Exec{DefaultExec("lang_nodejs")}, nil
}

var _ Detect = DetectNodejs // ensure interface is implemented

// readPackageJSON reads the package.json provided at input jsonpath. It returns any error encountered.
func readPackageJSON(jsonpath string) (PackageJSON, error) {
	bytes, err := os.ReadFile(jsonpath)
	if err != nil {
		return PackageJSON{}, fmt.Errorf("read file: %w", err)
	}

	var pkg PackageJSON
	if err := json.Unmarshal(bytes, &pkg); err != nil {
		return PackageJSON{}, fmt.Errorf("unmarshal: %w", err)
	}
	return pkg, pkg.Validate()
}

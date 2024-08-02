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
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/logger"
)

var _packageManagerRegexp = regexp.MustCompile(`^(npm|pnpm|yarn|bun)@\d+\.\d+\.\d+(-.+)?$`)

// PackageJSON represents the node package json descriptor.
type PackageJSON struct {
	Author         *string           `json:"author,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Files          []string          `json:"files,omitempty"`
	Keywords       []string          `json:"keywords,omitempty"`
	License        *string           `json:"license,omitempty"`
	Main           *string           `json:"main,omitempty"`
	Module         string            `json:"module,omitempty"`
	Name           string            `json:"name,omitempty"           validate:"required"`
	PackageManager string            `json:"packageManager,omitempty"`
	Private        bool              `json:"private,omitempty"`
	Scripts        map[string]string `json:"scripts,omitempty"`
	Version        string            `json:"version,omitempty"`

	PackageManagerName    string `json:"-"`
	PackageManagerVersion string `json:"-"`
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
func DetectNodejs(_ context.Context, log logger.Logger, destdir string, metadata Metadata) (Metadata, []Exec) {
	jsonpath := filepath.Join(destdir, craft.PackageJSON)
	pkg, err := readPackageJSON(jsonpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Warnf("failed to parse %s file: %s", jsonpath, err.Error())
		}
		return metadata, nil
	}

	log.Infof("nodejs detected, a '%s' is present and valid", craft.PackageJSON)

	// affect package manager and it's version separately for template generation
	if pkg.PackageManager == "" {
		pkg.PackageManager = "pnpm"
	}
	name, version, _ := strings.Cut(pkg.PackageManager, "@")
	if name != "" {
		pkg.PackageManagerName = name
	}
	if version != "" {
		pkg.PackageManagerVersion = version
	}

	metadata.Languages["nodejs"] = pkg
	metadata.ProjectName = pkg.Name
	if pkg.Main != nil {
		metadata.Binaries++
	}

	// deactivate makefile because commands are facilitated by package.json scripts
	if !metadata.NoMakefile {
		log.Warn("makefile option not available with nodejs generation, deactivating it")
		metadata.NoMakefile = true
	}

	return metadata, []Exec{DefaultExec("lang_nodejs")}
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

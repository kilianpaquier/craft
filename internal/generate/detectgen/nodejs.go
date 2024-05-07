package detectgen

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
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

//go:generate go-builder-generator generate -f nodejs.go -s PackageJSON -d tests

// PackageJSON represents the node package json descriptor.
type PackageJSON struct {
	Author                    *string           `json:"author,omitempty"`
	Description               *string           `json:"description,omitempty"`
	Files                     []string          `json:"files,omitempty"`
	Keywords                  []string          `json:"keywords,omitempty"`
	License                   *string           `json:"license,omitempty"`
	Main                      *string           `json:"main,omitempty"`
	Module                    string            `json:"module,omitempty"`
	Name                      string            `json:"name,omitempty"           validate:"required"`
	PackageManagerWithVersion string            `json:"packageManager,omitempty"`
	Private                   bool              `json:"private,omitempty"`
	Scripts                   map[string]string `json:"scripts,omitempty"`
	Version                   string            `json:"version,omitempty"`

	PackageManagerName    string `json:"-"`
	PackageManagerVersion string `json:"-"`
}

// Validate validates the given PackageJSON struct.
func (p *PackageJSON) Validate() error {
	var errs []error

	packageManager := regexp.MustCompile(`^(npm|pnpm|yarn|bun)@\d+\.\d+\.\d+(-.+)?$`)
	if p.PackageManagerWithVersion != "" && !packageManager.MatchString(p.PackageManagerWithVersion) {
		// json schema takes care of saying which regexp must be validated
		errs = append(errs, errors.New("package.json packageManager isn't valid"))
	}

	if err := validator.New().Struct(p); err != nil {
		errs = append(errs, fmt.Errorf("struct validation: %w", err))
	}
	return errors.Join(errs...)
}

// detectNodejs handles nodejs detection at configuration provided destination directory.
// It scans the project for a package.json and validates it.
//
// It returns the slice of GenerateFunc related to nodejs.
func detectNodejs(ctx context.Context, config *models.GenerateConfig) []GenerateFunc {
	log := logrus.WithContext(ctx)

	jsonpath := filepath.Join(config.Options.DestinationDir, models.PackageJSON)
	pkg, err := readPackageJSON(jsonpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.WithError(err).Warnf("failed to parse %s file", jsonpath)
		}
		return nil
	}

	log.Infof("nodejs detected, a %s is present and valid", models.PackageJSON)

	// affect package manager and it's version separately for template generation
	if pkg.PackageManagerWithVersion == "" {
		pkg.PackageManagerWithVersion = "pnpm"
	}
	name, version, _ := strings.Cut(pkg.PackageManagerWithVersion, "@")
	if name != "" {
		pkg.PackageManagerName = name
	}
	if version != "" {
		pkg.PackageManagerVersion = version
	}

	config.Languages[string(NameNodejs)] = pkg
	config.ProjectName = pkg.Name
	// automatically add one binary because there's no such things
	if pkg.Main != nil {
		config.Binaries++
	}

	// deactivate makefile because commands are facilitated by package.json scripts
	if !config.NoMakefile {
		log.Warn("makefile option not available with nodejs generation, deactivating it")
		config.NoMakefile = true
	}

	return []GenerateFunc{GetGenerateFunc(NameNodejs)}
}

// readPackageJSON reads the package.json provided at input jsonpath. It returns any error encountered.
func readPackageJSON(jsonpath string) (*PackageJSON, error) {
	bytes, err := os.ReadFile(jsonpath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var pkg *PackageJSON
	if err := json.Unmarshal(bytes, &pkg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return pkg, pkg.Validate()
}

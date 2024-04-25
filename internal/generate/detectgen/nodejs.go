package detectgen

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

//go:generate go-builder-generator generate -f nodejs.go -s PackageJSON -d builders

// PackageJSON represents the node package json descriptor.
type PackageJSON struct {
	Author         *string `json:"author,omitempty"`
	Description    *string `json:"description,omitempty"`
	License        *string `json:"license,omitempty"`
	Main           *string `json:"main,omitempty"`
	Name           string  `json:"name,omitempty"           validate:"required"`
	PackageManager *string `json:"packageManager,omitempty"`
	Private        bool    `json:"private,omitempty"`
	Version        string  `json:"version,omitempty"`
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

	config.Languages[string(NameNodejs)] = pkg
	config.ProjectName = pkg.Name
	// automatically add one binary because there's no such things
	if pkg.Main != nil {
		config.Binaries++
	}

	// automatically set default package manager if none was given
	if config.PackageManager == nil {
		config.PackageManager = lo.ToPtr("pnpm")
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

	if err := validator.New().Struct(pkg); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	return pkg, nil
}

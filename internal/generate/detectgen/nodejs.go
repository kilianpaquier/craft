package detectgen

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

// packageJSON represents the node package json descriptor.
type packageJSON struct {
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

	packagejson := filepath.Join(config.Options.DestinationDir, models.PackageJSON)
	bytes, err := os.ReadFile(packagejson)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.WithError(err).Info("failed to read package.json")
		}
		return nil
	}

	var descriptor packageJSON
	if err := json.Unmarshal(bytes, &descriptor); err != nil {
		log.WithError(err).Info("failed to unmarshal package.json")
		return nil
	}

	if err := validator.New().Struct(descriptor); err != nil {
		log.WithError(err).Error("invalid package.json file, proceeding without nodejs generation")
		return nil
	}

	log.Infof("nodejs detected, a %s is present and valid", models.PackageJSON)

	config.Languages = append(config.Languages, string(NameNodejs))
	config.ProjectName = descriptor.Name

	// automatically set default package manager if none was given
	if config.PackageManager == nil {
		config.PackageManager = lo.ToPtr("pnpm")
	}

	// deactivate makefile because commands are facilitated by package.json scripts
	if !config.NoMakefile {
		log.Warn("makefile option not available with nodejs generation, deactivating it")
		config.NoMakefile = true
	}

	// automatically add one binary because there's no such things
	if descriptor.Main != nil {
		config.Binaries++
	}
	return []GenerateFunc{GetGenerateFunc(NameNodejs)}
}

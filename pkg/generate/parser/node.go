package parser

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

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// FilePackageJSON represents package.json filename.
const FilePackageJSON = "package.json"

var (
	// ErrInvalidStaticDeployment represents the returned error when static deployment is provided in craft configuration
	// but current node project doesn't have a main file provided in package.json.
	ErrInvalidStaticDeployment = errors.New("package.json 'main' isn't provided but a static deployment is configured")

	// ErrInvalidPackageManager is the error returned when packageManager is missing or invalid in package.json.
	ErrInvalidPackageManager = errors.New("package.json 'packageManager' is missing or isn't valid")
)

var packageManagerRegexp = regexp.MustCompile(`^(npm|pnpm|yarn|bun)@\d+\.\d+\.\d+(-.+)?$`)

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

	if p.PackageManager != "" && !packageManagerRegexp.MatchString(p.PackageManager) {
		// json schema takes care of saying which regexp must be validated
		errs = append(errs, ErrInvalidPackageManager)
	}

	if err := validator.New().Struct(p); err != nil {
		errs = append(errs, fmt.Errorf("struct validation: %w", err))
	}
	return errors.Join(errs...)
}

// Node handles node repository parsing at destdir.
//
// It scans the project for a package.json and validates it.
func Node(ctx context.Context, destdir string, config *craft.Config) error {
	jsonpath := filepath.Join(destdir, FilePackageJSON)
	file, err := readPackageJSON(jsonpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("read package.json: %w", err)
		}
		return nil
	}
	generate.GetLogger(ctx).Infof("node detected, a '%s' is present and valid", FilePackageJSON)

	config.SetLanguage("node", file)
	config.ProjectName = file.Name
	if file.Main != nil {
		config.SetWorker("main")
	}
	return nil
}

var _ generate.Parser[craft.Config] = Node // ensure interface is implemented

// readPackageJSON reads package.json provided at input jsonpath.
// It returns any error encountered.
func readPackageJSON(jsonpath string) (PackageJSON, error) {
	bytes, err := os.ReadFile(jsonpath)
	if err != nil {
		return PackageJSON{}, fmt.Errorf("read file: %w", err)
	}

	var file PackageJSON
	if err := json.Unmarshal(bytes, &file); err != nil {
		return PackageJSON{}, fmt.Errorf("unmarshal: %w", err)
	}
	return file, file.Validate()
}

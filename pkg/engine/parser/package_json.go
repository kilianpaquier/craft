package parser

import (
	"errors"
	"regexp"
)

// FilePackageJSON represents package.json filename.
const FilePackageJSON = "package.json"

var (
	// ErrMissingPackageName is the error returned when package.json name is missing.
	ErrMissingPackageName = errors.New("package.json 'name' is missing")

	// ErrInvalidPackageManager is the error returned when packageManager is missing or invalid in package.json.
	ErrInvalidPackageManager = errors.New("package.json 'packageManager' is missing or isn't valid")
)

var packageManagerRegexp = regexp.MustCompile(`^(npm|pnpm|yarn|bun)@\d+\.\d+\.\d+(-.+)?$`)

// PackageJSON represents the node package json file.
type PackageJSON struct {
	Author         *string  `json:"author,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Files          []string `json:"files,omitempty"`
	Keywords       []string `json:"keywords,omitempty"`
	License        *string  `json:"license,omitempty"`
	Main           *string  `json:"main,omitempty"`
	Module         string   `json:"module,omitempty"`
	Name           string   `json:"name,omitempty"`
	PackageManager string   `json:"packageManager,omitempty"`
	Private        bool     `json:"private,omitempty"`
	PublishConfig  struct {
		Access     string `json:"access,omitempty"`
		Provenance bool   `json:"provenance,omitempty"`
		Registry   string `json:"registry,omitempty"`
		Tag        string `json:"tag,omitempty"`
	} `json:"publishConfig,omitempty"`
	Repository *struct {
		URL string `json:"url,omitempty"`
	} `json:"repository,omitempty"`
	Scripts map[string]string `json:"scripts,omitempty"`
	Version string            `json:"version,omitempty"`
}

// Validate validates the given PackageJSON struct.
func (p *PackageJSON) Validate() error {
	var errs []error

	if p.Name == "" {
		errs = append(errs, ErrMissingPackageName)
	}
	if !packageManagerRegexp.MatchString(p.PackageManager) {
		// json schema takes care of saying which regexp must be validated
		errs = append(errs, ErrInvalidPackageManager)
	}

	return errors.Join(errs...)
}

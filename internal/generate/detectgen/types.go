package detectgen

import (
	"context"

	filesystem "github.com/kilianpaquier/filesystem/pkg"

	"github.com/kilianpaquier/craft/internal/models"
)

// GenerateName represents a string alias for a generate reserved name.
type GenerateName string

const (
	// NameGeneric is the reserved name for generic template folder.
	NameGeneric GenerateName = "generic"
	// NameGolang is the reserved name for golang template folder.
	NameGolang GenerateName = "golang"
	// NameHelm is the reserved name for helm chart template folder.
	NameHelm GenerateName = "helm"
	// NameHugo is the reserved name for hugo template folder.
	NameHugo GenerateName = "hugo"
	// NameLicense is the reserved name for license template folder.
	NameLicense GenerateName = "license"
	// NameNodejs is the reserved name for nodejs template folder.
	NameNodejs GenerateName = "nodejs"
)

// ReservedNames returns the slice of string representing reserved folders name.
func ReservedNames() []string {
	return []string{
		string(NameGeneric),
		string(NameGolang),
		string(NameHelm),
		string(NameHugo),
		string(NameLicense),
		string(NameNodejs),
	}
}

// GenerateFunc is the signature function to implement to add a new language of framework templatization in craft.
type GenerateFunc func(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error

// DetectFunc is the signature function to implement to add a new language or framework detection in craft.
//
// The input configuration can be altered in any way.
type DetectFunc func(ctx context.Context, config *models.GenerateConfig) []GenerateFunc

// AllDetectFuncs returns the slice of all detects options,
// each one returning a slice of GenerateFunc in case the detection is truthy.
func AllDetectFuncs() []DetectFunc {
	return []DetectFunc{
		detectGolang,
		detectHelm,
		detectLicense,
		detectNodejs,
	}
}

package generate

import (
	"context"
)

// DetectGeneric represents the detection for generic projects (those without any associated implemented language).
//
// It returns the input metadata but modified with appropriate properties
// alongside the slice of ExecFunc to be executed to templatize the project.
func DetectGeneric(_ context.Context, _ string, metadata *Metadata) ([]ExecFunc, error) {
	if len(metadata.Languages) != 0 {
		return nil, nil
	}

	log.Warnf("no language detected, fallback to generic generation")
	if metadata.CI != nil {
		metadata.CI.Options = nil
	}
	return []ExecFunc{BasicExecFunc("lang_generic")}, nil
}

var _ DetectFunc = DetectGeneric // ensure interface is implemented

package generate

import (
	"context"

	"github.com/kilianpaquier/cli-sdk/pkg/clog"
)

// DetectGeneric represents the detection for generic projects (those without any associated implemented language).
//
// It returns the input metadata but modified with appropriate properties
// alongside the slice of Exec to be executed to templatize the project.
func DetectGeneric(_ context.Context, log clog.Logger, _ string, metadata *Metadata) ([]Exec, error) {
	if len(metadata.Languages) != 0 {
		return nil, nil
	}

	log.Warnf("no language detected, fallback to generic generation")
	if metadata.CI != nil {
		metadata.CI.Options = nil
	}
	return []Exec{DefaultExec("lang_generic")}, nil
}

var _ Detect = DetectGeneric // ensure interface is implemented

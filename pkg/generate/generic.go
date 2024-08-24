package generate

import (
	"context"

	"github.com/kilianpaquier/cli-sdk/pkg/clog"
)

// DetectGeneric represents the detection for generic projects (those without any associated implemented language).
//
// It returns the input metadata but modified with appropriate properties
// alongside the slice of Exec to be executed to templatize the project.
func DetectGeneric(_ context.Context, _ clog.Logger, _ string, metadata Metadata) (Metadata, []Exec, error) {
	if metadata.CI != nil {
		metadata.CI.Options = nil
	}
	return metadata, []Exec{DefaultExec("lang_generic")}, nil
}

var _ Detect = DetectGeneric // ensure interface is implemented

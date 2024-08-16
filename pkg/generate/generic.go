package generate

import (
	"context"
	"slices"

	"github.com/kilianpaquier/cli-sdk/pkg/clog"
	"github.com/samber/lo"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// DetectGeneric represents the detection for generic projects (those without any associated implemented language).
//
// It returns the input metadata but modified with appropriate properties
// alongside the slice of Exec to be executed to templatize the project.
func DetectGeneric(_ context.Context, _ clog.Logger, _ string, metadata Metadata) (Metadata, []Exec, error) {
	if metadata.CI != nil {
		// only keep generic applicable options
		options := lo.Filter(metadata.CI.Options, func(option string, _ int) bool {
			return slices.Contains([]string{craft.Dependabot, craft.Renovate}, option)
		})
		metadata.CI.Options = options
	}
	return metadata, []Exec{DefaultExec("lang_generic")}, nil
}

var _ Detect = DetectGeneric // ensure interface is implemented

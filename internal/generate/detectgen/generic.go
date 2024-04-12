package detectgen

import (
	"context"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/models"
)

// GenericFunc represents the detection for generic projects (those without any associated implemented language).
//
// It returns one slice element to generic templates from generic template folder.
func GenericFunc(ctx context.Context, config *models.GenerateConfig) []GenerateFunc {
	logrus.WithContext(ctx).Warn("no language detected, fallback to generic generation")

	if config.CI != nil {
		// only keep generic applicable options
		options := lo.Filter(config.CI.Options, func(option string, _ int) bool {
			return option == models.Dependabot || option == models.Renovate
		})
		config.CI.Options = options
	}
	return []GenerateFunc{GetGenerateFunc(NameGeneric)}
}

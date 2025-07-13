package generate

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/craft/generate/templates"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/generator"
)

// GeneratorGitignore downloads and writes .gitignore file in its right path.
//
// It patches it alongside with custom craft patches as some exclusion
// may be missing depending on craft layout generation.
func GeneratorGitignore(httpClient *http.Client) func(ctx context.Context, destdir string, config craft.Config) error {
	if httpClient == nil {
		httpClient = http.DefaultClient //nolint:revive
	}
	return func(ctx context.Context, destdir string, config craft.Config) error {
		mapping := map[string][]string{
			"go":    {"go"},
			"helm":  {"helm"},
			"hugo":  {"hugo"},
			"node":  {"node"},
			"shell": nil,
		}

		query := make([]string, 0, len(config.Languages)+3)
		for lang := range config.Languages {
			s, ok := mapping[lang]
			if ok {
				query = append(query, s...)
			}
		}
		query = append(query, "dotenv")

		if config.CI != nil {
			if slices.Contains(config.CI.Options, craft.Sonar) {
				query = append(query, "sonar", "sonarqube")
			}
		}

		if err := generator.DownloadGitignore(ctx, httpClient, filepath.Join(destdir, generator.FileGitignore), query...); err != nil {
			return fmt.Errorf("download gitignore: %w", err)
		}

		template := engine.Template[craft.Config]{
			Delimiters: engine.DelimitersBracket(),
			Patches:    []string{".gitignore" + engine.PatchExtension + engine.TmplExtension},
			Out:        ".gitignore",
		}
		if err := engine.ApplyTemplate(templates.FS(), destdir, template, config); err != nil {
			return fmt.Errorf("apply template: %w", err)
		}
		return nil
	}
}

var _ engine.Generator[craft.Config] = GeneratorGitignore(nil) // ensure interface is implemented

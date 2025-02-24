package templates

import (
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Misc returns the slice of templates globally related to a code repository (README.md, CODEOWNERS, etc.).
func Misc() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{"CODEOWNERS" + engine.TmplExtension},
			Out:        "CODEOWNERS",
		},
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{"README.md" + engine.TmplExtension},
			Out:        "README.md",
		},
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{".pre-commit-config.yaml" + engine.TmplExtension},
			Out:        ".pre-commit-config.yaml",
			Remove:     func(config craft.Config) bool { return slices.Contains(config.Exclude, craft.PreCommit) },
		},
	}
}

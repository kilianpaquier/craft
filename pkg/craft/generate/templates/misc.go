package templates

import (
	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Codeowners returns the slice of templates related to code owners configuration.
func Codeowners() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{"CODEOWNERS" + engine.TmplExtension},
			Out:        "CODEOWNERS",
		},
	}
}

// Readme returns the slice of templates related to README.md generation.
func Readme() []engine.Template[craft.Config] {
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{"README.md" + engine.TmplExtension},
			Out:        "README.md",
		},
	}
}

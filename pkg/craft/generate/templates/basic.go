package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Git is the handler for git specific files generation.
func Git() []engine.Template[craft.Config] {
	name := ".gitignore"
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      engine.Globs(name),
			Out:        name,
		},
	}
}

// Makefile is the handler for Makefile(s) generation.
func Makefile() []engine.Template[craft.Config] {
	var templates []engine.Template[craft.Config] //nolint:prealloc

	makefile := "Makefile"
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      []string{makefile + engine.TmplExtension},
		Out:        makefile,
		Remove: func(config craft.Config) bool {
			_, ok := config.Languages["node"] // don't generate makefiles with node
			return config.NoMakefile || ok
		},
	})

	scripts := path.Join("scripts", "mk")
	for _, src := range []string{path.Join(scripts, "build.mk")} {
		templates = append(templates, engine.Template[craft.Config]{
			Delimiters: engine.DelimitersBracket(),
			Globs:      engine.Globs(src),
			Out:        src,
			Remove: func(config craft.Config) bool {
				_, ok := config.Languages["node"] // don't generate makefiles with node
				return config.NoMakefile || ok
			},
		})
	}

	craftmk := path.Join(scripts, "craft.mk")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      []string{craftmk + engine.TmplExtension},
		Out:        craftmk,
		Remove: func(config craft.Config) bool {
			_, ok := config.Languages["node"] // don't generate makefiles with node
			return config.NoMakefile || ok
		},
	})

	return templates
}

// Readme is the handler for README.md generation.
func Readme() []engine.Template[craft.Config] {
	name := "README.md"
	return []engine.Template[craft.Config]{
		{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{name + engine.TmplExtension},
			Out:        name,
		},
	}
}

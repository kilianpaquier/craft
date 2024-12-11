package handler

import (
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// Defaults returns the full slice of handlers implemented in handler package.
func Defaults(handlers ...generate.Handler) []generate.Handler {
	return slices.Concat(
		[]generate.Handler{
			CodeCov,
			Dependabot,
			Docker,
			Git,
			GitHub,
			GitLab,
			Golang,
			Helm,
			Makefile,
			Readme,
			Renovate,
			SemanticRelease,
			Sonar,
		},

		// append custom handlers
		handlers,
	)
}

// PartGlob returns the glob string for HandlerResult globs parts.
//
// It should be used when templating a file split into multiple templates.
//
// The result is of the form: Dockerfile-*.part.tmpl, ci-*.part.tmpl, .gitlab-ci-*.part.yml .gitignore-*.part.tmpl, etc.
func PartGlob(src, name string) string {
	n := strings.TrimSuffix(name, filepath.Ext(name))
	if n == "" {
		n = name
	}
	glob := fmt.Sprint(n, "-*", craft.PartExtension, craft.TmplExtension)
	return path.Join(path.Dir(src), glob)
}

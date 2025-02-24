package templates

import (
	"path"
	"slices"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Chart returns the slice of templates related to helm chart generation.
func Chart() []engine.Template[craft.Config] {
	var templates []engine.Template[craft.Config] //nolint:prealloc

	tmplfiles := []string{
		path.Join("chart", "templates", "_helpers.tpl"),
		path.Join("chart", "templates", "configmap.yaml"),
		path.Join("chart", "templates", "cronjob.yaml"),
		path.Join("chart", "templates", "deployment.yaml"),
		path.Join("chart", "templates", "hpa.yaml"),
		path.Join("chart", "templates", "job.yaml"),
		path.Join("chart", "templates", "service.yaml"),
		path.Join("chart", "templates", "serviceaccount.yaml"),
	}
	for _, src := range tmplfiles {
		templates = append(templates, engine.Template[craft.Config]{
			Delimiters: engine.DelimitersChevron(),
			Globs:      []string{src + engine.TmplExtension},
			Out:        src,
			Remove:     func(config craft.Config) bool { return slices.Contains(config.Exclude, craft.Chart) },
		})
	}

	chartfiles := []string{
		path.Join("chart", ".craft"),
		path.Join("chart", ".helmignore"),
		path.Join("chart", "Chart.yaml"),
		path.Join("chart", "charts", ".gitkeep"),
	}
	for _, src := range chartfiles {
		templates = append(templates, engine.Template[craft.Config]{
			Delimiters: engine.DelimitersBracket(),
			Globs:      []string{src + engine.TmplExtension},
			Out:        src,
			Remove:     func(config craft.Config) bool { return slices.Contains(config.Exclude, craft.Chart) },
		})
	}

	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      engine.GlobsWithPart(path.Join("chart", "values.yaml")),
		Out:        path.Join("chart", "values.yaml"),
		Remove:     func(config craft.Config) bool { return slices.Contains(config.Exclude, craft.Chart) },
	})

	return templates
}

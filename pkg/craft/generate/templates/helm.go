package templates

import (
	"path"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
)

// Helm is the handler for chart folder generation.
func Helm() []engine.Template[craft.Config] {
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
			Remove:     func(config craft.Config) bool { return config.NoChart },
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
			Remove:     func(config craft.Config) bool { return config.NoChart },
		})
	}

	values := path.Join("chart", "values.yaml")
	templates = append(templates, engine.Template[craft.Config]{
		Delimiters: engine.DelimitersBracket(),
		Globs:      engine.Globs(values),
		Out:        values,
		Remove:     func(config craft.Config) bool { return config.NoChart },
	})

	return templates
}

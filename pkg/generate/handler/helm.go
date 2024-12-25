package handler

import (
	"path"
	"strings"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// Helm is the handler for chart folder generation.
func Helm(src, dest, name string) (generate.HandlerResult[craft.Config], bool) {
	handlers := []generate.Handler[craft.Config]{
		// files related to dir chart/templates
		helmTemplates,
		// files related to dir chart/charts
		helmCharts,
		// files related to dir chart
		helmConfig,
	}
	for _, handler := range handlers {
		if result, ok := handler(src, dest, name); ok {
			return result, ok
		}
	}
	return generate.HandlerResult[craft.Config]{}, false
}

func helmTemplates(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	// files related to dir chart/templates
	if !strings.Contains(src, path.Join("chart", "templates", name)) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterChevron(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return config.NoChart },
	}
	return result, true
}

func helmCharts(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	// files related to dir chart/charts
	if !strings.Contains(src, path.Join("chart", "charts", name)) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterChevron(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return config.NoChart },
	}
	return result, true
}

func helmConfig(src, _, name string) (generate.HandlerResult[craft.Config], bool) {
	// files related to dir chart
	if !strings.Contains(src, path.Join("chart", name)) {
		return generate.HandlerResult[craft.Config]{}, false
	}

	result := generate.HandlerResult[craft.Config]{
		Delimiter:    generate.DelimiterBracket(),
		Globs:        []string{src},
		ShouldRemove: func(config craft.Config) bool { return config.NoChart },
	}

	if name == "values.yaml" {
		result.Globs = append(result.Globs, PartGlob(src, name))
	}
	return result, true
}

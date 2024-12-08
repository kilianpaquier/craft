package handler

import (
	"path"
	"strings"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"

	"github.com/kilianpaquier/craft/pkg/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// Helm is the handler for chart folder generation.
func Helm(src, dest, name string) (generate.HandlerResult, bool) {
	handlers := []generate.Handler{
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
	return generate.HandlerResult{}, false
}

func helmTemplates(src, dest, name string) (generate.HandlerResult, bool) {
	// files related to dir chart/templates
	if !strings.Contains(src, path.Join("chart", "templates", name)) {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterChevron(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return metadata.NoChart },
	}
	return result, true
}

func helmCharts(src, dest, name string) (generate.HandlerResult, bool) {
	// files related to dir chart/charts
	if !strings.Contains(src, path.Join("chart", "charts", name)) {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterChevron(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return metadata.NoChart },
	}
	return result, true
}

func helmConfig(src, dest, name string) (generate.HandlerResult, bool) {
	// files related to dir chart
	if !strings.Contains(src, path.Join("chart", name)) {
		return generate.HandlerResult{}, false
	}

	result := generate.HandlerResult{
		Delimiter:      generate.DelimiterBracket(),
		Globs:          []string{src},
		ShouldGenerate: func(generate.Metadata) bool { return IsGenerated(dest) },
		ShouldRemove:   func(metadata generate.Metadata) bool { return metadata.NoChart },
	}

	switch name {
	case craft.File:
		result.ShouldGenerate = func(generate.Metadata) bool { return !cfs.Exists(dest) }
	case "values.yaml":
		result.Globs = append(result.Globs, PartGlob(src, name))
	}
	return result, true
}

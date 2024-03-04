package templating

import (
	"strings"
	"text/template"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// FuncMap returns a minimal template.FuncMap.
//
// Available functions are fromYaml, map, prefix, suffix and toYaml.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"fromYaml": fromYAML,
		"map":      mergeMaps,
		"toYaml":   toYAML,
	}
}

func mergeMaps(dst map[string]any, src ...any) map[string]any {
	for _, in := range src {
		var cast map[string]any
		if err := mapstructure.Decode(in, &cast); err != nil {
			dst["Error"] = err.Error()
			continue
		}
		if err := mergo.Merge(&dst, cast); err != nil {
			dst["Error"] = err.Error()
			continue
		}
	}
	return dst
}

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
// Copy of https://github.com/helm/helm/blob/main/pkg/engine/funcs.go.
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

// fromYAML converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
// Copy of https://github.com/helm/helm/blob/main/pkg/engine/funcs.go.
func fromYAML(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

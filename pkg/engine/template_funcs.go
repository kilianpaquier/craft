package engine

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"github.com/go-viper/mapstructure/v2"
	"github.com/goccy/go-yaml"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

// FuncMap returns a minimal template.FuncMap.
//
// It can be extended with MergeMaps.
func FuncMap(root string) template.FuncMap {
	return template.FuncMap{
		"cutAfter": cutAfter,
		"map":      mergeMaps,
		"glob":     glob(root),
		"toQuery":  toQuery,
		"toYaml":   toYAML,
	}
}

// glob returns a function checking with files.Glob if the input glob is present in root or its subdirectories.
func glob(root string) func(glob string) []string {
	return func(glob string) []string {
		return files.Glob(root, glob)
	}
}

// cutAfter cuts the input string at the first separator appearance
// and returns the resulting string.
func cutAfter(in, sep string) string {
	out, _, _ := strings.Cut(in, sep)
	return out
}

// mergeMaps mergs all src maps (an error is added to result map if those aren't maps) into dst map.
func mergeMaps(dst map[string]any, src ...any) map[string]any {
	for i, in := range src {
		var cast map[string]any
		if err := mapstructure.Decode(in, &cast); err != nil {
			dst[strconv.Itoa(i)+"_decode_error"] = err.Error()
			continue
		}
		if err := mergo.Merge(&dst, cast); err != nil {
			dst[strconv.Itoa(i)+"_merge_error"] = err.Error()
			continue
		}
	}
	return dst
}

// toQuery transforms a specific into its query parameter format.
func toQuery(in string) string {
	return url.QueryEscape(in)
}

// toYAML takes an interface, marshals it to yaml, and returns a string.
// It will always return a string, even on marshal error (empty string).
//
// This is designed to be called from a go template.
func toYAML(v any) string {
	b, err := yaml.MarshalWithOptions(v, yaml.Indent(2))
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(bytes.TrimSuffix(b, []byte("\n")))
}

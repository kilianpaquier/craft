package templating

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"github.com/go-viper/mapstructure/v2"
	"github.com/goccy/go-yaml"

	"github.com/kilianpaquier/craft/internal/helpers"
)

// FuncMap returns a minimal template.FuncMap.
//
// It can be extended with MergeMaps.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"cutAfter": CutAfter,
		"fromPtr":  helpers.FromPtr[string],
		"map":      MergeMaps,
		"toQuery":  ToQuery,
		"toYaml":   ToYAML,
	}
}

// CutAfter cuts the input string at the first separator appearance
// and returns the resulting string.
func CutAfter(in, sep string) string {
	out, _, _ := strings.Cut(in, sep)
	return out
}

// MergeMaps mergs all src maps (an error is added to result map if those aren't maps) into dst map.
func MergeMaps(dst map[string]any, src ...any) map[string]any {
	for i, in := range src {
		var cast map[string]any
		if err := mapstructure.Decode(in, &cast); err != nil {
			dst[fmt.Sprint(i, "_decode_error")] = err.Error()
			continue
		}
		if err := mergo.Merge(&dst, cast); err != nil {
			dst[fmt.Sprint(i, "_merge_error")] = err.Error()
			continue
		}
	}
	return dst
}

// ToQuery transforms a specific into its query parameter format.
func ToQuery(in string) string {
	return url.QueryEscape(in)
}

// ToYAML takes an interface, marshals it to yaml, and returns a string.
// It will always return a string, even on marshal error (empty string).
//
// This is designed to be called from a go template.
func ToYAML(v any) string {
	b, err := yaml.MarshalWithOptions(v, yaml.Indent(2))
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(bytes.TrimSuffix(b, []byte("\n")))
}

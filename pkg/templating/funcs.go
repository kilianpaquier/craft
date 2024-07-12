package templating

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/go-viper/mapstructure/v2"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

// FuncMap returns a minimal template.FuncMap.
//
// It can be extended with MergeMaps.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"map":     MergeMaps,
		"toQuery": ToQuery,
		"toYaml":  ToYAML,
	}
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
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(v); err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

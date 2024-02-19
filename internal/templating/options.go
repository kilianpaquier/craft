package templating

import (
	"text/template"
)

// Option represents a function taking as input a template to potentially alter it.
//
// Use with caution. You may use predefined options like WithDelims to avoid unexpected template modifications.
type Option func(*template.Template) *template.Template

// WithDelims is an option for ApplyTemplate to use custom go template delimiters.
func WithDelims(left, right string) Option {
	return func(t *template.Template) *template.Template {
		return t.Delims(left, right)
	}
}

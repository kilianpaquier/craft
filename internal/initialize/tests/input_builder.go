package tests

import (
	"strconv"
	"strings"

	"github.com/samber/lo"

	"github.com/kilianpaquier/craft/internal/models"
)

// InputBuilder is a builder to set the inputs during init command.
type InputBuilder struct {
	API            bool
	Chart          bool
	Maintainer     models.Maintainer
	OpenAPIVersion string
}

// NewInputBuilder creates a new InputBuilder.
func NewInputBuilder() *InputBuilder {
	return &InputBuilder{}
}

// SetMaintainer sets the maintainers.
func (r *InputBuilder) SetMaintainer(maintainer models.Maintainer) *InputBuilder {
	r.Maintainer = maintainer
	return r
}

// SetAPI sets the API response.
func (r *InputBuilder) SetAPI(api bool) *InputBuilder {
	r.API = api
	return r
}

// SetOpenAPIVersion sets the openapi version.
func (r *InputBuilder) SetOpenAPIVersion(openAPIVersion string) *InputBuilder {
	r.OpenAPIVersion = openAPIVersion
	return r
}

// SetChart sets the chart response.
func (r *InputBuilder) SetChart(chart bool) *InputBuilder {
	r.Chart = chart
	return r
}

// Build builds the InputBuilder into a string reader to use on init command.
func (r *InputBuilder) Build() (*strings.Reader, error) {
	values := []string{
		r.Maintainer.Name, "\n",
		lo.FromPtr(r.Maintainer.Email), "\n",
		lo.FromPtr(r.Maintainer.URL), "\n",
		strconv.FormatBool(r.API), "\n",
		r.OpenAPIVersion, "\n",
		strconv.FormatBool(r.Chart), "\n",
	}
	return strings.NewReader(strings.Join(values, "")), nil
}

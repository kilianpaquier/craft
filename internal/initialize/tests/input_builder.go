package tests

import (
	"errors"
	"strconv"
	"strings"

	"github.com/samber/lo"

	"github.com/kilianpaquier/craft/internal/models"
	models_tests "github.com/kilianpaquier/craft/internal/models/tests"
)

// InputBuilder is a builder to set the inputs during init command.
type InputBuilder struct {
	*models_tests.CraftConfigBuilder
	API   string
	Chart string
}

// NewInputBuilder creates a new InputBuilder.
func NewInputBuilder() *InputBuilder {
	return &InputBuilder{
		CraftConfigBuilder: models_tests.NewCraftConfigBuilder(),
	}
}

// SetDescription sets the description.
func (r *InputBuilder) SetDescription(description string) *InputBuilder {
	r.CraftConfigBuilder.SetDescription(description)
	return r
}

// SetMaintainers sets the maintainers.
func (r *InputBuilder) SetMaintainers(maintainer ...models.Maintainer) *InputBuilder {
	r.CraftConfigBuilder.SetMaintainers(maintainer...)
	return r
}

// SetAPI sets the API response.
func (r *InputBuilder) SetAPI(api string) *InputBuilder {
	r.API = api
	b, _ := strconv.ParseBool(api)
	r.SetNoAPI(!b)
	return r
}

// SetOpenAPIVersion sets the openapi version.
func (r *InputBuilder) SetOpenAPIVersion(openAPIVersion string) *InputBuilder {
	r.CraftConfigBuilder.SetOpenAPIVersion(openAPIVersion)
	return r
}

// SetChart sets the chart response.
func (r *InputBuilder) SetChart(chart string) *InputBuilder {
	r.Chart = chart
	b, _ := strconv.ParseBool(chart)
	r.SetNoChart(!b)
	return r
}

// Build builds the InputBuilder into a string reader to use on init command.
func (r *InputBuilder) Build() (*strings.Reader, error) {
	if len(r.Maintainers) == 0 {
		return nil, errors.New("maintainers are mandatory")
	}

	values := []string{
		lo.FromPtr(r.Description), "\n",
		r.Maintainers[0].Name, "\n",
		lo.FromPtr(r.Maintainers[0].Email), "\n",
		lo.FromPtr(r.Maintainers[0].URL), "\n",
		r.API, "\n",
		r.OpenAPIVersion, "\n",
		r.Chart, "\n",
	}
	return strings.NewReader(strings.Join(values, "")), nil
}

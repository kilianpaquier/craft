package models

//go:generate go-builder-generator generate -f craft.go -s CraftConfig,Maintainer,GenerateConfig,GenerateOptions,CI,API,Docker -d tests

//go:generate fieldalignment -fix craft.go

//go:generate tagalign -fix -sort -strict -order "json,yaml,builder,validate" .

const (
	// CraftFile is the craft configuration file name.
	CraftFile = ".craft"

	// TmplExtension is the extension for templates file.
	TmplExtension = ".tmpl"

	// Github is just the github constant.
	Github = "github"
	// Gitlab is just the gitlab constant.
	Gitlab = "gitlab"
	// GoCmd represents the cmd folder where go main.go should be placed according to go layout.
	GoCmd = "cmd"
	// GoMod represents the go.mod filename.
	GoMod = "go.mod"
	// License represents the target filename for the generated project LICENSE.
	License = "LICENSE"
	// SwaggerFile represents the target filename for the generated project api.yml.
	SwaggerFile = "api.yml"
)

const (
	// CodeCov is the codecov option for CI tuning.
	CodeCov string = "codecov"
	// Dependabot is the dependabot option for CI tuning.
	Dependabot string = "dependabot"
	// Sonar is the sonar option for CI tuning.
	Sonar string = "sonar"
)

// CraftConfig represents all options configurable in .craft file at root project.
//
// yaml tags are for .craft file and json tags for templating.
type CraftConfig struct {
	API          *API         `json:"api,omitempty"         yaml:"api,omitempty"                            validate:"omitempty,required"`
	CI           *CI          `json:"-"                     yaml:"ci,omitempty"                             validate:"omitempty,required"`
	Description  *string      `json:"description,omitempty" yaml:"description,omitempty"`
	Docker       *Docker      `json:"docker,omitempty"      yaml:"docker,omitempty"                         validate:"omitempty,required"`
	License      *string      `json:"-"                     yaml:"license,omitempty"                        validate:"omitempty,oneof=agpl-3.0 apache-2.0 bsd-2-clause bsd-3-clause bsl-1.0 cc0-1.0 epl-2.0 gpl-2.0 gpl-3.0 lgpl-2.1 mit mpl-2.0 unlicense"`
	Maintainers  []Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"   builder:"append" validate:"required,dive,required"`
	NoChart      bool         `json:"-"                     yaml:"no_chart,omitempty"`
	NoGoreleaser bool         `json:"-"                     yaml:"no_goreleaser,omitempty"`
	NoMakefile   bool         `json:"-"                     yaml:"no_makefile,omitempty"`
}

// CI is the struct for craft ci tuning.
type CI struct {
	Name    string   `json:"-" yaml:"name,omitempty"                     validate:"required,oneof=github gitlab"`
	Options []string `json:"-" yaml:"options,omitempty" builder:"append" validate:"omitempty,dive,oneof=codecov dependabot sonar"`
}

// API is the struct for craft api tuning.
type API struct {
	OpenAPIVersion *string `json:"-" yaml:"openapi_version,omitempty" validate:"omitempty,oneof=v2 v3"`
}

// Docker is the struct for craft docker tuning.
type Docker struct {
	Registry *string `json:"registry,omitempty" yaml:"registry,omitempty"`
	Port     *uint16 `json:"port,omitempty"     yaml:"port,omitempty"`
}

// Maintainer represents a project maintainer. It's inspired from helm Maintainer struct.
//
// The only difference are the present tags and the pointers on both email and url properties.
type Maintainer struct {
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`
	URL   *string `json:"url,omitempty"   yaml:"url,omitempty"`
	Name  string  `json:"name,omitempty"  yaml:"name,omitempty"`
}

// GenerateConfig represents an extension of CraftConfig during craft generate command.
//
// Those additional properties will be enriched during generate execution and project parsing.
// They will be used for files and helm chart templating (if applicable).
//
// yaml tags are for .craft file and json tags for templating.
type GenerateConfig struct {
	CraftConfig
	Clis          map[string]struct{} `json:"-"                     yaml:"-"`
	Crons         map[string]struct{} `json:"crons,omitempty"       yaml:"-"`
	Jobs          map[string]struct{} `json:"jobs,omitempty"        yaml:"-"`
	Workers       map[string]struct{} `json:"workers,omitempty"     yaml:"-"`
	ModuleName    string              `json:"-"                     yaml:"-"`
	ModuleVersion string              `json:"-"                     yaml:"-"`
	ProjectName   string              `json:"projectName,omitempty" yaml:"-"`

	Options GenerateOptions `json:"-" yaml:"-"`
}

// GenerateOptions represents all options configurable in craft generate command.
type GenerateOptions struct {
	EndDelim   string `validate:"required"`
	StartDelim string `validate:"required"`

	DestinationDir string `validate:"required"`
	TemplatesDir   string

	Force    []string `builder:"append" validate:"omitempty,dive,required"`
	ForceAll bool
}

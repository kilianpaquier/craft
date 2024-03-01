package models

//go:generate fieldalignment -fix .

//go:generate tagalign -fix -sort -strict -order "json,yaml,builder,validate" .

//go:generate go-builder-generator generate -f craft.go -s CraftConfig,Maintainer,GenerateConfig,GenerateOptions -d tests

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

// CraftConfig represents all options configurable in .craft file at root project.
//
// yaml tags are for .craft file and json tags for templating.
type CraftConfig struct {
	Description    *string      `json:"description,omitempty"    yaml:"description,omitempty"`
	DockerRegistry *string      `json:"dockerRegistry,omitempty" yaml:"docker_registry,omitempty"`
	License        *string      `json:"-"                        yaml:"license,omitempty"                          validate:"omitempty,oneof=agpl-3.0 apache-2.0 bsd-2-clause bsd-3-clause bsl-1.0 cc0-1.0 epl-2.0 gpl-2.0 gpl-3.0 lgpl-2.1 mit mpl-2.0 unlicense"`
	Port           *uint16      `json:"port,omitempty"           yaml:"port,omitempty"`
	OpenAPIVersion string       `json:"-"                        yaml:"openapi_version,omitempty"                  validate:"omitempty,oneof=v2 v3"`
	CI             string       `json:"-"                        yaml:"ci,omitempty"                               validate:"omitempty,oneof=gitlab github"`
	Maintainers    []Maintainer `json:"maintainers,omitempty"    yaml:"maintainers,omitempty"     builder:"append" validate:"required,dive,required"`
	CodeCov        bool         `json:"-"                        yaml:"codecov,omitempty"`
	NoAPI          bool         `json:"noAPI,omitempty"          yaml:"no_api,omitempty"`
	NoChart        bool         `json:"-"                        yaml:"no_chart,omitempty"`
	NoDockerfile   bool         `json:"-"                        yaml:"no_dockerfile,omitempty"`
	NoGoreleaser   bool         `json:"-"                        yaml:"no_goreleaser,omitempty"`
	NoMakefile     bool         `json:"-"                        yaml:"no_makefile,omitempty"`
	Sonar          bool         `json:"-"                        yaml:"sonar,omitempty"`
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
	Clis        map[string]struct{} `json:"-"                     yaml:"-"`
	Crons       map[string]struct{} `json:"crons,omitempty"       yaml:"-"`
	Jobs        map[string]struct{} `json:"jobs,omitempty"        yaml:"-"`
	Workers     map[string]struct{} `json:"workers,omitempty"     yaml:"-"`
	ModuleName  string              `json:"-"                     yaml:"-"`
	ProjectName string              `json:"projectName,omitempty" yaml:"-"`

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

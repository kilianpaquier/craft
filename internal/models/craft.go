package models

//go:generate go-builder-generator generate -f craft.go -s CraftConfig,Maintainer,GenerateConfig,GenerateOptions,CI,Docker,Release -d tests

const (
	// CraftFile is the craft configuration file name.
	CraftFile = ".craft"

	// TmplExtension is the extension for templates file.
	TmplExtension = ".tmpl"

	// Gocmd represents the cmd folder where go main.go should be placed according to go layout.
	Gocmd = "cmd"
	// Gomod represents the go.mod filename.
	Gomod = "go.mod"
	// PackageJSON represents the package.json filename.
	PackageJSON = "package.json"

	// License represents the target filename for the generated project LICENSE.
	License = "LICENSE"
)

const (
	// Bitbucket is just the bitbucket constant.
	Bitbucket = "bitbucket"
	// Gitea is just the gitea constant.
	Gitea = "gitea"
	// Github is just the github constant.
	Github = "github"
	// Gitlab is just the gitlab constant.
	Gitlab = "gitlab"
)

const (
	// CodeCov is the codecov option for CI tuning.
	CodeCov string = "codecov"
	// CodeQL is the codeql option for CI tuning.
	CodeQL string = "codeql"
	// Dependabot is the dependabot option for CI tuning.
	Dependabot string = "dependabot"
	// Pages is the pages option for CI tuning.
	Pages string = "pages"
	// Renovate is the renovate option for CI tuning.
	Renovate string = "renovate"
	// Sonar is the sonar option for CI tuning.
	Sonar string = "sonar"
)

const (
	// GithubApps is the value for github release mode with a github app.
	GithubApps string = "github-apps"
	// GithubToken is the value for github release mode with a github token.
	GithubToken string = "github-token"
	// PersonalToken is the value for github release mode with a personal token (PAT).
	PersonalToken string = "personal-token"
)

// AllOptions returns the slice with all availables CI options.
func AllOptions() []string {
	return []string{CodeCov, CodeQL, Dependabot, Pages, Renovate, Sonar}
}

// CraftConfig represents all options configurable in .craft file at root project.
//
// yaml tags are for .craft file and json tags for templating.
type CraftConfig struct {
	CI           *CI          `json:"-"                     yaml:"ci,omitempty"                             validate:"omitempty,required"`
	Description  *string      `json:"description,omitempty" yaml:"description,omitempty"`
	Docker       *Docker      `json:"docker,omitempty"      yaml:"docker,omitempty"                         validate:"omitempty,required"`
	License      *string      `json:"-"                     yaml:"license,omitempty"                        validate:"omitempty,oneof=agpl-3.0 apache-2.0 bsd-2-clause bsd-3-clause bsl-1.0 cc0-1.0 epl-2.0 gpl-2.0 gpl-3.0 lgpl-2.1 mit mpl-2.0 unlicense"`
	Maintainers  []Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"   builder:"append" validate:"required,dive,required"`
	NoChart      bool         `json:"-"                     yaml:"no_chart,omitempty"`
	NoGoreleaser bool         `json:"-"                     yaml:"no_goreleaser,omitempty"`
	NoMakefile   bool         `json:"-"                     yaml:"no_makefile,omitempty"`
	Platform     string       `json:"-"                     yaml:"platform,omitempty"                       validate:"omitempty,oneof=bitbucket gitea github gitlab"`
}

// CI is the struct for craft continuous integration tuning.
type CI struct {
	Name    string   `json:"-" yaml:"name,omitempty"                     validate:"required,oneof=github gitlab"`
	Release Release  `json:"-" yaml:"release,omitempty"                  validate:"required"`
	Options []string `json:"-" yaml:"options,omitempty" builder:"append" validate:"omitempty,dive,oneof=codecov codeql dependabot pages renovate sonar"`
}

// Release is the struct for craft continuous integration release specifics configuration.
type Release struct {
	Auto      bool   `json:"-" yaml:"auto"`
	Backmerge bool   `json:"-" yaml:"backmerge"`
	Disable   bool   `json:"-" yaml:"disable,omitempty"`
	Mode      string `json:"-" yaml:"mode,omitempty"   validate:"omitempty,oneof=github-apps personal-token github-token"`
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
// yaml tags are for .craft file and json tags for helm templating.
type GenerateConfig struct {
	CraftConfig

	Binaries    uint8               `json:"-"                     yaml:"-"`
	Clis        map[string]struct{} `json:"-"                     yaml:"-"`
	Crons       map[string]struct{} `json:"crons,omitempty"       yaml:"-"`
	Jobs        map[string]struct{} `json:"jobs,omitempty"        yaml:"-"`
	Languages   map[string]any      `json:"-"                     yaml:"-" builder:"append"`
	ProjectHost string              `json:"projectHost"           yaml:"-"`
	ProjectName string              `json:"projectName,omitempty" yaml:"-"`
	ProjectPath string              `json:"projectPath"           yaml:"-"`
	Workers     map[string]struct{} `json:"workers,omitempty"     yaml:"-"`

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

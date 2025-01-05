package craft

import (
	"slices"

	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// Config represents all options configurable in .craft file at root project.
//
// Note that yaml tags are for .craft file property keys
// and json tags for templating data.
type Config struct {
	parser.Executables `yaml:",inline"`

	// Bot represents the name of the maintenance bot (renovate, dependabot, etc).
	//
	// It's optional and some restrictions may apply (see craft JSON schema).
	// For instance, when working with GitLab, only Renovate is supported.
	Bot *string `json:"-" yaml:"bot,omitempty"`

	// CI is the structure containing all optional and configurable properties for CI purposes.
	CI *CI `json:"-" yaml:"ci,omitempty"`

	// Description represents the project description.
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`

	// Docker is the structure containing all optional and configurable properties
	// for Dockerfile (and Helm in case a chart is generated).
	Docker *Docker `json:"docker,omitempty" yaml:"docker,omitempty"`

	// License is the project license name.
	License *string `json:"-" yaml:"license,omitempty"`

	// Languages is a map of languages name with its specificities.
	Languages map[string]any `json:"-" yaml:"-"`

	// Maintainers is the slice of all project maintainers.
	Maintainers []*Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"`

	// NoChart can be given to avoid generating an Helm chart.
	//
	// By default, an Helm chart is generated since a project could satisfies one of the following possibilities:
	//	- the project is just an Helm chart for another product
	// 	- the project is a product with an Helm chart (with one or multiple resources, cronjobs, job, worker, etc)
	NoChart bool `json:"-" yaml:"no_chart,omitempty"`

	// NoGoreleaser can be given to avoid generating a .goreleaser.yml file.
	//
	// By default, if a given project is a Go project,
	// and a cmd CLI is defined (cmd/<some useful CLI name>)
	// a .goreleaser.yml file is generated.
	//
	// As such, it's unnecessary to specify this property when your project isn't a Go one.
	NoGoreleaser bool `json:"-" yaml:"no_goreleaser,omitempty"`

	// NoMakefile can be given to avoid generating a Makefile and additional Makefiles in scripts/mk/*.mk.
	//
	// When crafting a Node project, it's unnecessary to specify this property since no Makefile will be generated anyway.
	// It's because Node projects contain all their scripts in package.json.
	NoMakefile bool `json:"-" yaml:"no_makefile,omitempty"`

	parser.VCS `yaml:",inline"` // put at the end to get sorted properties (Platform especially) in written yaml file.
}

// Auth contains all authentication methods related to CI configuration.
type Auth struct {
	// Maintenance represents the authentication method for the maintenance bot (renovate, dependabot, etc.).
	Maintenance *string `json:"-" yaml:"maintenance,omitempty"`

	// Release represents the authentication method for the release process (GitHub Token, Personal Access Token, etc.).
	// It's unavailable when working with a GitLab project.
	Release *string `json:"-" yaml:"release,omitempty"`
}

// CI is the struct for craft continuous integration tuning.
type CI struct {
	// Auth contains all authentication methods related to CI configuration.
	Auth Auth `json:"-" yaml:"auth,omitempty"`

	// Name represents the CI name (GitHub, GitLab, etc.).
	//
	// Note that those must be in lowercase.
	Name string `json:"-" yaml:"name,omitempty"`

	// Options is the slice of CI options.
	Options []string `json:"-" yaml:"options,omitempty"`

	// Release is the struct containing all tuning around release process (auto release, backmerge, etc.).
	Release *Release `json:"-" yaml:"release,omitempty"`

	// Static is the struct containing all tuning around static deployment (auto, name, etc.).
	Static *Static `json:"-" yaml:"static,omitempty"`
}

// Docker is the struct for craft docker tuning.
type Docker struct {
	// Port represents the port to expose in the Dockerfile / Helm chart.
	//
	// It's shared for all cmd executables that could be defined in the project.
	Port *uint16 `json:"port,omitempty" yaml:"port,omitempty"`

	// Registry represents the Docker registry to use.
	Registry *string `json:"registry,omitempty" yaml:"registry,omitempty"`
}

// Maintainer represents a project maintainer. It's inspired from helm Maintainer struct.
//
// The only difference are the present tags and the pointers on both email and url properties.
type Maintainer struct {
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`
	Name  string  `json:"name,omitempty"  yaml:"name,omitempty"`
	URL   *string `json:"url,omitempty"   yaml:"url,omitempty"`
}

// Release is the struct for craft continuous integration release specifics configuration.
type Release struct {
	// Auto is the boolean indicating whether release should be done automatically on default branch.
	Auto bool `json:"-" yaml:"auto,omitempty"`

	// Backmerge is the boolean indicating whether backmerge should be done during release process (with semantic-release).
	Backmerge bool `json:"-" yaml:"backmerge,omitempty"`
}

// Static represents the configuration for static deployment.
type Static struct {
	// Auto is the boolean indicating whether static deployment
	// should be done automatically on default branch.
	Auto bool `json:"-" yaml:"auto,omitempty"`

	// Name is the name of the static deployment (netlify, pages, etc.).
	Name string `json:"-" yaml:"name,omitempty"`
}

// IsBot returns truthy in case the input bot is the one specified in configuration.
//
// It returns false if no maintenance bot is specified in configuration.
func (c Config) IsBot(bot string) bool {
	return c.Bot != nil && *c.Bot == bot
}

// IsCI returns truthy in case the input name is the one specified in configuration.
//
// It returns false if CI is disabled.
func (c Config) IsCI(name string) bool {
	return c.CI != nil && c.CI.Name == name
}

// HasDockerRegistry returns truthy in case the configuration has a docker registry configuration.
func (c Config) HasDockerRegistry() bool {
	return c.Docker != nil && c.Docker.Registry != nil
}

// IsMaintenanceAuth returns truthy in case the input auth value is the one specified in configuration maintenance auth.
//
// It returns false if neither CI or auth maintenance isn't specified in configuration.
func (c Config) IsMaintenanceAuth(auth string) bool {
	return c.CI != nil && c.CI.Auth.Maintenance != nil && *c.CI.Auth.Maintenance == auth
}

// IsReleaseAuth returns truthy in case the input auth value is the one specified in configuration release auth.
//
// It returns false if neither CI or auth release isn't specified in configuration.
func (c Config) IsReleaseAuth(auth string) bool {
	return c.CI != nil && c.CI.Auth.Release != nil && *c.CI.Auth.Release == auth
}

// HasRelease returns truthy in case the configuration has CI enabled and Release configuration.
func (c Config) HasRelease() bool {
	return c.CI != nil && c.CI.Release != nil
}

// IsAutoRelease returns truthy in case the configuration has CI enabled, release enabled and auto actived.
func (c Config) IsAutoRelease() bool {
	return c.CI != nil && c.CI.Release != nil && c.CI.Release.Auto
}

// SetLanguage sets a language with its specificities.
func (c *Config) SetLanguage(name string, value any) {
	if c.Languages == nil {
		c.Languages = map[string]any{}
	}
	c.Languages[name] = value
}

// IsStatic returns truthy in case the input static value is the one specified in configuration as static name.
//
// It returns false in case there's no CI or no Static configuration.
func (c Config) IsStatic(static string) bool {
	return c.CI != nil && c.CI.Static != nil && c.CI.Static.Name == static
}

// EnsureDefaults migrates old properties into new fields and ensures default properties are always sets.
func (c *Config) EnsureDefaults() {
	c.retroCompatibility()

	// small sanitization for CI configuration part
	func() {
		if c.CI == nil {
			return
		}
		slices.Sort(c.CI.Options)
	}()
}

func (*Config) retroCompatibility() {
	// TBD in case a migration is needed
}

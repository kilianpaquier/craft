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
	parser.Executables `yaml:"-"`

	// Bot represents the name of the maintenance bot (renovate, dependabot, etc).
	//
	// It's optional and some restrictions may apply (see craft JSON schema).
	// For instance, when working with GitLab, only Renovate is supported.
	Bot string `json:"-" yaml:"bot,omitempty"`

	// CI is the structure containing all optional and configurable properties for CI purposes.
	CI *CI `json:"ci,omitempty" yaml:"ci,omitempty"`

	// Description represents the project description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Exclude is the slice of string indicating which part of generation must not be made.
	//
	// JSON schema can be followed to get more information on that part
	// and which parts can be excluded.
	Exclude []string `json:"-" yaml:"exclude,omitempty"`

	// Include is the slice of string indicating optional part to generate (additionally to base generation).
	//
	// JSON schema can be followed to get more information on that part
	// and which parts can be included.
	Include []string `json:"-" yaml:"include,omitempty"`

	// Languages is a map of languages name with its specificities.
	Languages map[string]any `json:"-" yaml:"-"`

	// License is the project license name.
	License string `json:"-" yaml:"license,omitempty"`

	// Maintainers is the slice of all project maintainers.
	Maintainers []*Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"`

	parser.VCS `yaml:",inline"` // put at the end to get sorted properties (Platform especially) in written yaml file.
}

// Auth contains all authentication methods related to CI configuration.
type Auth struct {
	// Maintenance represents the authentication method for the maintenance bot (renovate, dependabot, etc.).
	Maintenance string `json:"-" yaml:"maintenance,omitempty"`

	// Release represents the authentication method for the release process (GitHub Token, Personal Access Token, etc.).
	// It's unavailable when working with a GitLab project.
	Release string `json:"-" yaml:"release,omitempty"`
}

// CI is the struct for craft continuous integration tuning.
type CI struct {
	// Auth contains all authentication methods related to CI configuration.
	Auth Auth `json:"-" yaml:"auth,omitempty"`

	// Deployment is the struct containing all tuning around deployment.
	Deployment *Deployment `json:"-" yaml:"deployment,omitempty"`

	// Docker is the structure containing all optional and configurable properties
	// for Dockerfile (and Helm docker properties like the registry).
	Docker *Docker `json:"docker,omitempty" yaml:"docker,omitempty"`

	// Helm is the structure containing all optional and configurable properties
	// for Helm
	Helm *Helm `json:"-" yaml:"helm,omitempty"`

	// Name represents the CI name (GitHub, GitLab, etc.).
	//
	// Note that those must be in lowercase.
	Name string `json:"-" yaml:"name,omitempty"`

	// Options is the slice of CI options.
	Options []string `json:"-" yaml:"options,omitempty"`

	// Release is the struct containing all tuning around release process (auto release, backmerge, etc.).
	Release *Release `json:"-" yaml:"release,omitempty"`
}

// Helm represents the configuration for Helm.
type Helm struct {
	// Path is the chart fullname.
	// It's optional and by default will be computed with VCS informations (owner/repository).
	//
	// Example:
	//  - https://github.com/kilianpaquier/craft will give kilianpaquier/craft
	//
	// To ensure the packaged chart is pushed to a specific registry, please use Registry property.
	Path string `json:"-" yaml:"path,omitempty"`

	// Publish is an enum string composed of 'manual', 'auto', 'none' to indicate
	// whether the Helm chart should be publish on an helm repository.
	Publish string `json:"-" yaml:"publish,omitempty"`

	// Registry represents the Helm registry to use.
	//
	// When Registry is given, push action will be generated in Continuous Integration.
	Registry string `json:"-" yaml:"registry,omitempty"`
}

// Deployment represents the configuration for deployment.
type Deployment struct {
	// Auto is the boolean indicating whether deployment
	// should be done automatically on default branch.
	Auto bool `json:"-" yaml:"auto,omitempty"`

	// Platform is the deployment platform name (netlify, pages, kubernetes, azure, gcp, etc.).
	Platform string `json:"-" yaml:"platform,omitempty"`
}

// Docker is the struct for craft docker tuning.
type Docker struct {
	// Path is the image fullname.
	// It's optional and by default will be computed with VCS informations (owner/repository).
	//
	// Example:
	//  - https://github.com/kilianpaquier/craft will give kilianpaquier/craft
	//
	// To ensure the docker image is pushed to a specific registry, please use Registry property.
	Path string `json:"path,omitempty" yaml:"path,omitempty"`

	// Port represents the port to expose in the Dockerfile / Helm chart.
	//
	// It's shared for all cmd executables that could be defined in the project.
	Port uint16 `json:"port,omitempty" yaml:"port,omitempty"`

	// Registry represents the Docker registry to use.
	Registry string `json:"registry,omitempty" yaml:"registry,omitempty"`
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

// IsCI returns truthy in case the input name is the one specified in configuration.
//
// It returns false if CI is disabled.
func (c Config) IsCI(name string) bool {
	return c.CI != nil && c.CI.Name == name
}

// IsMaintenanceAuth returns truthy in case the input auth value is the one specified in configuration maintenance auth.
//
// It returns false if neither CI or auth maintenance isn't specified in configuration.
func (c Config) IsMaintenanceAuth(auth string) bool {
	return c.CI != nil && c.CI.Auth.Maintenance == auth
}

// IsReleaseAuth returns truthy in case the input auth value is the one specified in configuration release auth.
//
// It returns false if neither CI or auth release isn't specified in configuration.
func (c Config) IsReleaseAuth(auth string) bool {
	return c.CI != nil && c.CI.Auth.Release == auth
}

// HasRelease returns truthy in case the configuration has CI enabled and Release configuration.
func (c Config) HasRelease() bool {
	return c.CI != nil && c.CI.Release != nil
}

// IsAutoRelease returns truthy in case the configuration has CI enabled, release enabled and auto actived.
func (c Config) IsAutoRelease() bool {
	return c.CI != nil && c.CI.Release != nil && c.CI.Release.Auto
}

// HasHelmPublish returns truthy in case the configuration has CI enabled, helm chart generation enabled
// and publication to an helm repository enabled.
func (c Config) HasHelmPublish() bool {
	return c.CI != nil && c.CI.Helm != nil && slices.Contains([]string{HelmAuto, HelmManual}, c.CI.Helm.Publish)
}

// IsDeployment returns truthy in case the input platform value is the one specified in configuration as deployment platform name.
//
// It returns false in case there's no CI or no Deployment configuration.
func (c Config) IsDeployment(platform string) bool {
	return c.CI != nil && c.CI.Deployment != nil && c.CI.Deployment.Platform == platform
}

// SetLanguage sets a language with its specificities.
func (c *Config) SetLanguage(name string, value any) {
	if c.Languages == nil {
		c.Languages = map[string]any{}
	}
	c.Languages[name] = value
}

// EnsureDefaults migrates old properties into new fields
// and ensures default properties are always sets.
func (c *Config) EnsureDefaults() {
	c.retroCompatibility()

	slices.Sort(c.Exclude)
	slices.Sort(c.Include)

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

package craft

import (
	"slices"
)

// Configuration represents all options configurable in .craft file at root project.
//
// Note that yaml tags are for .craft file property keys and json tags for templating data.
type Configuration struct {
	Bot          *string       `json:"-"                     yaml:"bot,omitempty"`
	CI           *CI           `json:"-"                     yaml:"ci,omitempty"`
	Description  *string       `json:"description,omitempty" yaml:"description,omitempty"`
	Docker       *Docker       `json:"docker,omitempty"      yaml:"docker,omitempty"`
	License      *string       `json:"-"                     yaml:"license,omitempty"`
	Maintainers  []*Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"`
	NoChart      bool          `json:"-"                     yaml:"no_chart,omitempty"`
	NoGoreleaser bool          `json:"-"                     yaml:"no_goreleaser,omitempty"`
	NoMakefile   bool          `json:"-"                     yaml:"no_makefile,omitempty"`
	NoReadme     bool          `json:"-"                     yaml:"no_readme,omitempty"`
	Platform     string        `json:"-"                     yaml:"platform,omitempty"`
}

// Auth contains all authentication methods related to CI configuration.
type Auth struct {
	Maintenance *string `json:"-" yaml:"maintenance,omitempty"`
	Release     *string `json:"-" yaml:"release,omitempty"`
}

// CI is the struct for craft continuous integration tuning.
type CI struct {
	Auth    Auth     `json:"-" yaml:"auth,omitempty"    validate:"omitempty,required"`
	Name    string   `json:"-" yaml:"name,omitempty"    validate:"required"`
	Options []string `json:"-" yaml:"options,omitempty"`
	Release *Release `json:"-" yaml:"release,omitempty" validate:"omitempty,required"`
	Static  *Static  `json:"-" yaml:"static,omitempty"  validate:"omitempty,required"`
}

// Docker is the struct for craft docker tuning.
type Docker struct {
	Port     *uint16 `json:"port,omitempty"     yaml:"port,omitempty"`
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
	Auto      bool `json:"-" yaml:"auto,omitempty"`
	Backmerge bool `json:"-" yaml:"backmerge,omitempty"`
}

// Static represents the configuration for static deployment.
type Static struct {
	Auto bool   `json:"-" yaml:"auto,omitempty"`
	Name string `json:"-" yaml:"name,omitempty"`
}

// IsBot returns truthy in case the input bot is the one specified in configuration.
//
// It returns false if no maintenance bot is specified in configuration.
func (c Configuration) IsBot(bot string) bool {
	return c.Bot != nil && *c.Bot == bot
}

// IsCI returns truthy in case the input name is the one specified in configuration.
//
// It returns false if CI is disabled.
func (c Configuration) IsCI(name string) bool {
	return c.CI != nil && c.CI.Name == name
}

// HasDockerRegistry returns truthy in case the configuration has a docker registry configuration.
func (c Configuration) HasDockerRegistry() bool {
	return c.Docker != nil && c.Docker.Registry != nil
}

// IsMaintenanceAuth returns truthy in case the input auth value is the one specified in configuration maintenance auth.
//
// It returns false if neither CI or auth maintenance isn't specified in configuration.
func (c Configuration) IsMaintenanceAuth(auth string) bool {
	return c.CI != nil && c.CI.Auth.Maintenance != nil && *c.CI.Auth.Maintenance == auth
}

// IsReleaseAuth returns truthy in case the input auth value is the one specified in configuration release auth.
//
// It returns false if neither CI or auth release isn't specified in configuration.
func (c Configuration) IsReleaseAuth(auth string) bool {
	return c.CI != nil && c.CI.Auth.Release != nil && *c.CI.Auth.Release == auth
}

// HasRelease returns truthy in case the configuration has CI enabled and Release configuration.
func (c Configuration) HasRelease() bool {
	return c.CI != nil && c.CI.Release != nil
}

// IsAutoRelease returns truthy in case the configuration has CI enabled, release enabled and auto actived.
func (c Configuration) IsAutoRelease() bool {
	return c.CI != nil && c.CI.Release != nil && c.CI.Release.Auto
}

// IsStatic returns truthy in case the input static value is the one specified in configuration as static name.
//
// It returns false in case there's no CI or no Static configuration.
func (c Configuration) IsStatic(static string) bool {
	return c.CI != nil && c.CI.Static != nil && c.CI.Static.Name == static
}

// EnsureDefaults migrates old properties into new fields and ensures default properties are always sets.
func (c *Configuration) EnsureDefaults() {
	c.retroCompatibility()

	// small sanitization for CI configuration part
	func() {
		if c.CI == nil {
			return
		}
		slices.Sort(c.CI.Options)
	}()
}

func (*Configuration) retroCompatibility() {
	// TBD in case a migration is needed
}

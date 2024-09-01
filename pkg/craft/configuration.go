package craft

// Configuration represents all options configurable in .craft file at root project.
//
// yaml tags are for .craft file and json tags for templating.
type Configuration struct {
	Bot          *string      `json:"-"                     yaml:"bot,omitempty"                            validate:"omitempty,oneof=dependabot renovate"`
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

// Auth contains all authentication methods related to CI configuration.
type Auth struct {
	Maintenance *string `json:"-" yaml:"maintenance,omitempty" validate:"omitempty,oneof=github-app github-token mend.io personal-token"`
	Release     *string `json:"-" yaml:"release,omitempty"     validate:"omitempty,oneof=github-app github-token personal-token"`
}

// CI is the struct for craft continuous integration tuning.
type CI struct {
	Auth    Auth     `json:"-" yaml:"auth,omitempty"                     validate:"omitempty,required"`
	Name    string   `json:"-" yaml:"name,omitempty"                     validate:"required"`
	Options []string `json:"-" yaml:"options,omitempty" builder:"append"`
	Release *Release `json:"-" yaml:"release"                            validate:"omitempty,required"`
	Static  *Static  `json:"-" yaml:"static,omitempty"                   validate:"omitempty,required"`
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
	Name  string  `json:"name,omitempty"  yaml:"name,omitempty"  validate:"required"`
	URL   *string `json:"url,omitempty"   yaml:"url,omitempty"`
}

// Release is the struct for craft continuous integration release specifics configuration.
type Release struct {
	Action    string `json:"-" yaml:"action"              validate:"required,oneof=gh-release release-drafter release-please semantic-release"`
	Auto      bool   `json:"-" yaml:"auto,omitempty"`
	Backmerge bool   `json:"-" yaml:"backmerge,omitempty"`
}

// Static represents the configuration for static deployment.
type Static struct {
	Auto bool   `json:"-" yaml:"auto,omitempty"`
	Name string `json:"-" yaml:"name,omitempty" validate:"required,oneof=netlify pages"`
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

// IsReleaseAction returns truthy in case the input action is the one specified by the configuration release action.
//
// It returns false if there's no CI or Release specified in configuration.
func (c Configuration) IsReleaseAction(action string) bool {
	return c.CI != nil && c.CI.Release != nil && c.CI.Release.Action == action
}

// IsStatic returns truthy in case the input static value is the one specified in configuration as static name.
//
// It returns false in case there's no CI or no Static configuration.
func (c Configuration) IsStatic(static string) bool {
	return c.CI != nil && c.CI.Static != nil && c.CI.Static.Name == static
}

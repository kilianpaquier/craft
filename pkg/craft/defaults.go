package craft

import (
	"slices"

	"github.com/kilianpaquier/craft/internal/helpers"
)

// EnsureDefaults acts to ensure default properties are always sets
// and migrates old properties into new fields.
func (c *Configuration) EnsureDefaults() {
	c.retroCompatibility()

	// ensure defaults values are set for maintenance bot
	if c.Bot != nil {
		if *c.Bot == Dependabot {
			c.CI.Auth.Maintenance = nil // dependabot doesn't need any mode
		}

		if c.Platform == Gitlab {
			c.Bot = helpers.ToPtr(Renovate) // dependabot is not available on craft for gitlab
			c.CI.Auth.Maintenance = nil     // renovate on gitlab isn't configurable
		}
	}

	c.ensureDefaultCI()
}

func (c *Configuration) ensureDefaultCI() {
	if c.CI == nil {
		return
	}
	slices.Sort(c.CI.Options)

	// ensure default values are set for CI
	// ...

	// specific gitlab CICD
	if c.CI.Name == Gitlab {
		c.CI.Options = slices.DeleteFunc(c.CI.Options, func(option string) bool { return option == Labeler }) // labeler isn't available on gitlab CICD
	}

	func() {
		if c.CI.Release == nil {
			c.CI.Auth.Release = nil
			return
		}

		// ensure default values are set for release
		// ...

		if c.CI.Auth.Release == nil {
			c.CI.Auth.Release = helpers.ToPtr(GithubToken) // set default release mode for github actions
		}

		// specific gitlab CICD
		if c.CI.Name == Gitlab {
			c.CI.Auth.Release = nil // release auth isn't available with gitlab CICD
		}
	}()
}

func (c *Configuration) retroCompatibility() {
	if c.CI != nil {
		// generic function to match an option included in a slice of options
		del := func(options ...string) func(option string) bool {
			return func(option string) bool {
				return slices.Contains(options, option)
			}
		}

		// migrate old renovate / dependabot option
		switch {
		case slices.Contains(c.CI.Options, Dependabot):
			c.Bot = helpers.ToPtr(Dependabot)
			c.CI.Options = slices.DeleteFunc(c.CI.Options, del(Dependabot))
		case slices.Contains(c.CI.Options, Renovate):
			c.Bot = helpers.ToPtr(Renovate)
			c.CI.Options = slices.DeleteFunc(c.CI.Options, del(Renovate))
		}

		// migrate old netlify / pages option
		switch {
		case slices.Contains(c.CI.Options, Netlify):
			c.CI.Static = &Static{Name: Netlify}
			c.CI.Options = slices.DeleteFunc(c.CI.Options, del(Netlify))
		case slices.Contains(c.CI.Options, Pages):
			c.CI.Static = &Static{Name: Pages}
			c.CI.Options = slices.DeleteFunc(c.CI.Options, del(Pages))
		}
	}
}

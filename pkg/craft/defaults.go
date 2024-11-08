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
		if c.Platform == GitLab {
			c.Bot = helpers.ToPtr(Renovate) // dependabot is not available on craft for gitlab
		}
	}

	c.ensureDefaultCI()
}

func (c *Configuration) ensureDefaultCI() {
	if c.CI == nil {
		return
	}
	slices.Sort(c.CI.Options)

	if c.Bot != nil {
		if *c.Bot == Dependabot || c.Platform == GitLab {
			c.CI.Auth.Maintenance = nil // dependabot and gitlab don't need any mode
		}
	}

	// specific gitlab CICD
	if c.CI.Name == GitLab {
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
			c.CI.Auth.Release = helpers.ToPtr(GitHubToken) // set default release mode for github actions
		}

		// specific gitlab CICD
		if c.CI.Name == GitLab {
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

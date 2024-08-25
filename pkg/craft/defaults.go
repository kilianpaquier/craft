package craft

import (
	"slices"

	"github.com/kilianpaquier/craft/internal/helpers"
)

// EnsureDefaults acts to ensure default properties are always sets
// and migrates old properties into new fields.
func (c *Configuration) EnsureDefaults() {
	c.retroCompatibility() // nolint:revive

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

		if c.CI.Release.Action == "" {
			c.CI.Release.Action = GhRelease // set default release action in case it's not provided
		}
		if c.CI.Auth.Release == nil {
			c.CI.Auth.Release = helpers.ToPtr(GithubToken) // set default release mode for github actions
		}

		// specific gitlab CICD
		if c.CI.Name == Gitlab {
			c.CI.Auth.Release = nil               // release auth isn't available with gitlab CICD
			c.CI.Release.Action = SemanticRelease // only semantic release is available on gitlab CICD
		}

		// specific github actions (to put inside its own condition when a third CI name is implemented)
		if c.CI.Release.Action == ReleaseDrafter || c.CI.Release.Action == GhRelease {
			c.CI.Release.Backmerge = false // gh-release or release-drafter don't handle backmerge
			if !slices.Contains(c.CI.Options, Labeler) {
				c.CI.Options = append(c.CI.Options, Labeler) // Labeler is mandatory for gh-release and release-drafter since those releaser are based on pull requests labels
			}
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

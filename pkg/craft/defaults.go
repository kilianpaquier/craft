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
			c.Bot = helpers.ToPtr(Renovate) // dependabot is not available on craft for GitLab
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
		} else if *c.Bot == Renovate && c.CI.Auth.Maintenance == nil {
			c.CI.Auth.Maintenance = helpers.ToPtr(GitHubToken)
		}
	}

	// labeler is only available on GitHub Actions
	if c.CI.Name != GitHub {
		c.CI.Options = slices.DeleteFunc(c.CI.Options, func(option string) bool { return option == Labeler })
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

		// specific GitLab CICD
		if c.CI.Name == GitLab {
			c.CI.Auth.Release = nil // release auth isn't available with GitLab CICD
		}
	}()
}

func (*Configuration) retroCompatibility() {
	// TBD in case a migration is needed
}

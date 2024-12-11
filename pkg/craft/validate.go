package craft

import (
	"slices"
)

// Validate acts as multiple purposes:
//   - Ensure default properties are always sets
//   - Migrate old properties into new fields
//   - Validate that the current configuration is valid with craft JSON schema
func (c *Configuration) Validate() {
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

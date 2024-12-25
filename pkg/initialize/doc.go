/*
The initialize package provides functions to create a new project.

The main function to be used is Run and it can be tuned with options (see documentation).

Example:

	type Config struct { License *string }

	func main() {
		ctx := context.Background()
		destdir, _ := os.Getwd()

		config, err := initialize.Run(ctx, destdir, initialize.WithFormGroups(License))
		// handle err
	}

	// ReadLicense returns the appropriate huh.Group for initialize.Run form groups.
	func License(config *Config) *huh.Group {
		return huh.NewGroup(huh.NewInput().
			Title("Would you like to specify a license ?").
			Validate(func(s string) error {
				if s != "" {
					config.License = &s
				}
				return nil
			}))
	}

Example with craft pre-defined form groups:

	// Note that group.Maintainer may be provided multiple times
	// in case multiple maintainers are needed.

	func main() {
		ctx := context.Background()
		destdir, _ := os.Getwd()

		config, err := initialize.Run(ctx, destdir,
			initialize.WithFormGroups(group.Maintainer, group.Chart))
		// handle err
	}
*/
package initialize

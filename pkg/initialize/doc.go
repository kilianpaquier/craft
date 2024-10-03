/*
The initialize package provides functions to create a new craft project.

The main function to be used is Run and it can be tuned with options (see documentation).

Example:

	func main() {
		ctx := context.Background()
		destdir, _ := os.Getwd()

		config, err := initialize.Run(ctx, destdir)
	}

Example with a custom FormGroup:

	func main() {
		ctx := context.Background()
		destdir, _ := os.Getwd()

		config, err := initialize.Run(ctx, destdir, initialize.WithFormGroups(initialize.ReadMaintainer, initialize.ReadChart, ReadLicense))
		if err != nil {
			// handle err
		}
	}

	// ReadLicense returns the appropriate huh.Group for initialize.Run form groups.
	func ReadLicense(config *craft.Configuration) *huh.Group {
		return huh.NewGroup(huh.NewInput().
			Title("Would you like to specify a license ?").
			Validate(func(s string) error {
				if s != "" {
					config.License = &s
				}
				return nil
			}))
	}
*/
package initialize

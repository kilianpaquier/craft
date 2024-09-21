/*
The initialize package provides functions to create a new craft project.

The main function to be used is Run and it can be tuned with options (see documentation).

Multiple Run option are available to further tune the initialization of a project with craft: WithLogger, WithReader, WithInputReader.
For further information about those options, please consult their specific documentation.

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
	}

	// ReadLicense reads an answer from stream reader (by default os.Stdin or the one provided by WithReader)
	// to set the license name (or none) of the generated project.
	func ReadLicense(_ logger.Logger, config craft.Configuration, ask initialize.Ask) craft.Configuration {
		config.License = ask("Which license would you like to use ?")
		return config
	}
*/
package initialize

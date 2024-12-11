/*
Package initialize exposes pre-defined groups for craft configuration initialization with engine.Initialize.

Example:

	func main() {
		config, err := engine.Initialize(ctx, engine.WithFormGroups(initialize.Maintainer, initialize.License))
		// handle err
	}
*/
package initialize

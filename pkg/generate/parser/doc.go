/*
Package parser provides a bunch of implementations for Parser function from generate package.

Either those can be used indivually to enrich metadata given in input of every parser or as a whole with Defaults function.

Example:

	// one specific parser
	func main() {
		handlers, err := parser.Git(ctx, "path/to/dir", metadata) // metadata is updated in the process
		// handle err
	}

	// as a whole
	func main() {
		for _, p := range parser.Defaults() {
			handlers, err := p(ctx, "path/to/dir", metadata) // metadata is updated in the process
			// handle err
		}
	}

	// fully used with generate.Run
	func main() {
		// config (craft.Configuration) is updated during the process
		// and returned updated at the end
		config, err := generate.Run(ctx, config, generate.WithParsers(parser.Defaults()...))
		// handle err
	}
*/
package parser

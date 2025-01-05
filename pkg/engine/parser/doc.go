/*
Package parser provides a bunch of functions to be wrapped with generate.Parser function signature.

Examples:

	type config struct { ... }

	func Parser(ctx context.Context, destdir string, c *config) error {
		// do something to parse the repository / configuration
		return nil
	}

	var _ generate.Parser[config] = Parser // ensure interface is implemented

	// single parser call
	func main() {
		var c config
		err := parser.Parser(ctx, "path/to/dir", &c)
		// handle err
	}

	// fully used with engine.Generate
	func main() {
		var c config
		err := engine.Generate(ctx, &c, generate.WithParsers(Parser, ...))
		// handle err
	}
*/
package parser

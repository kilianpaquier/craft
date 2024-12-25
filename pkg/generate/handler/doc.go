/*
Package handler provides a bunch of implementations of Handler type from generate package.

Either those can be used indivually with generate.Run or as a whole either from as a Parser
or with an anonymous function implementing Parser type.

Examples:

	// one specific handler
	func main() {
		codecov := handler.CodeCov(config)
		result, ok := codecov("path/to/src/file", "path/to/dest/file", "filename")
		// ok will be true if the input file is handled by the current handler (in this case codecov)

		// remove will be true in case the file must be removed (should not be generated)
		// it may be true when the input config doesn't say "hey, I want the codecov file in my project" (in codecov case)
		remove := result.ShouldRemove()

		// generate will be true in case the file must be generated
		// it may be false depending on handler, either because the file is not a generated one ("Code generate by ...")
		// or already exists ...
		generate := result.ShouldGenerate()
	}

	// as a whole (with defaults slice)
	func main() {
		for _, h := range handler.Defaults(config) {
			result, ok := h("path/to/src/file", "path/to/dest/file", "filename")
			// handle file correctly
		}
	}

	// fully used with generate.Run
	func main() {
		config, err := generate.Run(ctx, config, generate.WithHandlers(handler.Defaults()...))
		// handle err
	}
*/
package handler

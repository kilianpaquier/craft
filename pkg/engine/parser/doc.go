/*
Package parser provides a bunch of functions to be wrapped with generate.Parser function signature.

Examples:

	type config struct { ... }

	func ParserNode(ctx context.Context, destdir string, c *config) error {
		var jsonfile parser.PackageJSON
		jsonpath := filepath.Join(destdir, parser.FilePackageJSON)
		if err := files.ReadJSON(jsonpath, &jsonfile, os.ReadFile); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return fmt.Errorf("read json: %w", err)
		}
		engine.GetLogger().Infof("node detected, a '%s' is present and valid", parser.FilePackageJSON)

		if err := jsonfile.Validate(); err != nil {
			return fmt.Errorf("validate '%s': %w", parser.FilePackageJSON, err)
		}
		// do something with parsed jsonfile (e.g. update config since it's a pointer)
		return nil
	}

	var _ generate.Parser[config] = ParserNode // ensure interface is implemented

	// single parser call
	func main() {
		var c config
		err := parser.ParserNode(ctx, "path/to/dir", &c)
		// handle err
	}

	// fully used with engine.Generate
	func main() {
		destdir, _ := os.Getwd()

		var c config
		err := engine.Generate(ctx, destdir, c, []engine.Parser[config]{ParserNode})
		// handle err
	}
*/
package parser

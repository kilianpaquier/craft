/*
Package generate exposes pre-defined parsers for craft repositories parsing to use with engine.Generate.

Example:

	func main() {
		ctx := context.Background()
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, craft.File)

		// read configuration
		var config craft.Config
		if err := files.ReadYAML(dest, &config, os.ReadFile); err != nil {
			logger.Fatal(err)
		}
		config.EnsureDefaults()

		// run generation
		options := []engine.GenerateOption[craft.Config]{
			...
			engine.WithParsers(generate.ParserGit, generate.ParserLicense, generate.ParserGolang, generate.ParserNode, generate.ParserHelm),
			...
		}
		config, err := engine.Generate(ctx, config, options...)
		// handle err
	}
*/
package generate

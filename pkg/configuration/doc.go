/*
Package configuration provides various features to read, write
and validate configuration files (specifically YAML ones).

Validation feature is based on JSON Schema (https://json-schema.org/).

Example:

	type Config struct { ... }

	func main() {
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, "filename.ext")

		read := func() ([]byte, error) { return os.ReadFile(dest) }
		if err := configuration.Validate(dest, read); err != nil {
			logger.Fatal(err)
		}

		var config Config
		if err := configuration.ReadYAML(dest, &config); err != nil {
			logger.Fatal(err)
		}

		...

		opts := []yaml.EncodeOption{
			yaml.Indent(2),
			yaml.IndentSequence(true),
			yaml.WithComment(yaml.CommentMap{
				"$": []*yaml.Comment{
					yaml.HeadComment(" Some specific configuration file\n---"),
				},
			}),
		}
		if err := configuration.WriteYAML(dest, config, opts...); err != nil {
			logger.Fatal(err)
		}
	}
*/
package configuration

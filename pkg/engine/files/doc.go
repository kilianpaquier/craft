/*
Package files provides various features to read, write
and validate files with JSON schema.

Validation feature is based on JSON Schema (https://json-schema.org/)
and santhosh-tekuri/jsonschema (https://github.com/santhosh-tekuri/jsonschema).

Example:

	type config struct { ... }

	func main() {
		destdir, _ := os.Getwd()
		dest := filepath.Join(destdir, "filename.ext")

		if err := files.Validate(
			func(out any) error { return files.ReadJSON("path/to/schema.json", out, customfs.ReadFile) },
			func(out any) error { return files.ReadYAML(dest, out, os.ReadFile) },
		); err != nil {
			logger.Fatal(err)
		}

		var c config
		if err := configuration.ReadYAML(dest, &config); err != nil {
			logger.Fatal(err)
		}

		...

		if err := configuration.WriteYAML(dest, config, yaml.Indent(2), yaml.IndentSequence(true)); err != nil {
			logger.Fatal(err)
		}
	}
*/
package files

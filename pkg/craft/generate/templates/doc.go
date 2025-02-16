/*
Package templates provides a bunch of functions returning engine.Template struct.
Those functions can be used with either All function, returning all concatenated templates
or each one alone get specific templates.

Those functions can be used with engine.ApplyTemplate, engine.ApplyPatches or engine.GeneratorTemplates.

Note that engine.ApplyPatches is called by default in engine.ApplyTemplate.

Examples:

	func main() {
		destdir, _ := os.Getwd()

		var c craft.Config
		for _, template := range templates.Chart() {
			err := engine.ApplyTemplate(templates.FS(), destdir, template, c)
			// handle err
		}
	}

	func main() {
		destdir, _ := os.Getwd()

		var c craft.Config
		for _, template := range templates.Chart() {
			err := engine.ApplyPatches(templates.FS(), destdir, template, c)
			// handle err
		}
	}

	func main() {
		destdir, _ := os.Getwd()
		var c craft.Config

		f := engine.GeneratorTemplates(templates.FS(), templates.Chart())
		err := f(ctx, destdir, c)
		// handle err
	}
*/
package templates

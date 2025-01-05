/*
Package templates provides a bunch of functions returning engine.Template struct.
Those functions can be used with either All function, returning all concatenated templates
or each one alone get specific templates.

Those functions are only useful with engine.WithTemplates function, option of engine.Generate function:

	engine.WithTemplates(templates.FS(), templates.All())
	engine.WithTemplates(os.DirFS("path/to/templates"), templates.Helm())
	...
*/
package templates

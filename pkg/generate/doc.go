/*
The generate package provides multiple functions to extend craft generation.

The main function to be used is Run and it can be tuned with options (see documentation).

Multiple Run option are available to further tune the generation of a project with craft: WithMetaHandlers, WithDelimiters, WithDetects, WithDestination, WithForce, WithForceAll, WithLogger, WithTemplates.
For further information about those options, please consult their specific documentation.

The best options for generation tuning are however WithMetaHandlers, WithDetects, WithTemplates.
Those three options allows to override default craft templates, enrich generated files conditions and even add new languages parsing and generation.

Example:

	func main() {
		ctx := context.Background()
		config := craft.Configuration{} // may be read and saved with craft package

		generate.SetLogger(clog.Std()) // set stdlib log, clog.Slog() can also be used or even a custom implement of clog.Logger interface

		config, err := generate.Run(ctx, config,
			generate.WithDelimiters("<<", ">>"),
			generate.WithDestination(destdir),

			// Detects returns the default slice of DetectFuncs
			generate.WithDetectFuncs(generate.Detects()...),

			// MetaHandlers returns the default slice of MetaHandlers
			// which is a slice of funcs each taking as input Metadata and returning a func handling a specific file
			// i.e. the default ones are related to Docker, GitHub Actions, GitLab CI/CD, goreleaser, Makefile and Sonar
			generate.WithMetaHandlers(generate.MetaHandlers()...),

			generate.WithForce(force...),
			generate.WithForceAll(forceAll),

			// override the templates, by default here FS is the embedded fs of craft which default templates
			// another possibility is cfs.OS which takes an implementation reading the current filesystem
			//
			// the first input string it the folder path where the templates are located
			generate.WithTemplates("templates", generate.FS()),
		)
		if err != nil {
			// handle err
		}
	}
*/
package generate

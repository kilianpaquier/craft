/*
Package engine provides functions to create and generate a project layout.

It contains two main functions, Initialize and Generate which split project initialization and project generation in two parts.

Initialize example:

	type config struct { ... }

	func main() {
		ctx := context.Background()
		destdir, _ := os.Getwd()

		config, err := engine.Initialize(ctx, destdir, engine.WithFormGroups(License))
		// handle err
	}

	func License(c *config) *huh.Group {
		var license bool
		return huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to specify a license (optional) ?").
				Value(&license),

			huh.NewSelect[string]().
				Title("Which one ?").
				OptionsFunc(func() []huh.Option[string] {
					if !license {
						return nil
					}
					return huh.NewOptions(licenses...)
				}, &license).
				Validate(func(s string) error {
					if s != "" {
						config.License = &s
					}
					return nil
				}),
		)
	}

Generate example:

	type config struct { ... }

	func main() {
		ctx := context.Background()
		destdir, _ := os.Getwd()

		// run generation
		options := []engine.GenerateOption[config]{
			engine.WithDestination[config](destdir),
			engine.WithLogger[config](logger),
			engine.WithParsers(ParserGit),
			// WithTemplates takes a fs.FS, as such, anything implementing this interface can be given (like embed.FS)
			engine.WithTemplates(os.DirFS("path/to/templates"), Templates()),
		}
		config, err := engine.Generate(ctx, config, options...)
		// handle err
	}

	func ParserGit(ctx context.Context, destdir string, config *config) error {
		vcs, err := parser.Git(destdir)
		if err != nil {
			engine.GetLogger(ctx).Warnf("failed to retrieve git vcs configuration: %v", err)
			return nil // a repository may not be a git repository
		}
		engine.GetLogger(ctx).Infof("git repository detected")

		config.VCS = vcs
		return nil
	}

	func Templates() []engine.Templates[config] {
		name := ".gitignore"
		return []engine.Template[craft.Config]{
			{
				Delimiters: engine.DelimitersBracket(),
				Globs:      engine.Globs(name),
				Out:        name,
				// Remove can be given to remove a specific file in some specific case instead of generating it
				Remove: func (config) bool { return false },
				// GeneratePolicy can be given to tune generation, see the appropriate documentation
				GeneratePolicy: engine.PolicyAlways,
			},
		}
	}
*/
package engine

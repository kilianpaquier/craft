package parser

import (
	"slices"

	"github.com/kilianpaquier/craft/pkg/configuration/craft"
	"github.com/kilianpaquier/craft/pkg/generate"
)

// Defaults returns the full slice of handlers implemented in parser package.
func Defaults(parsers ...generate.Parser[craft.Config]) []generate.Parser[craft.Config] {
	return slices.Concat(
		[]generate.Parser[craft.Config]{
			// parse git repository first
			Git,

			License, // parse license configuration in configuration and generate it
			Golang,  // parse go.mod
			Node,    // parse package.json
		},

		// append custom parsers
		parsers,

		[]generate.Parser[craft.Config]{
			Helm, // parse helm configuration and overrides (must be call last)
		},
	)
}

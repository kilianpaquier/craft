package craft

import "github.com/goccy/go-yaml"

// EncodeOpts returns the options related to YAML encoding with goccy/go-yaml.
func EncodeOpts() []yaml.EncodeOption {
	return []yaml.EncodeOption{
		yaml.Indent(2),
		yaml.IndentSequence(true),
		yaml.WithComment(yaml.CommentMap{
			"$": []*yaml.Comment{
				yaml.HeadComment(
					" Craft configuration file (https://github.com/kilianpaquier/craft)",
					" yaml-language-server: $schema=https://raw.githubusercontent.com/kilianpaquier/craft/beta/.schemas/craft.schema.json",
				),
			},
		}),
	}
}

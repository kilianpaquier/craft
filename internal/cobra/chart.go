package cobra

import (
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/craft/pkg/craft/generate"
	"github.com/kilianpaquier/craft/pkg/craft/generate/templates"
	"github.com/kilianpaquier/craft/pkg/engine"
)

var chartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Generate project layout's helm chart",
	Run: gen(
		// keep most of parser to get a complete chart
		engine.WithParsers(generate.ParserGit, generate.ParserGolang, generate.ParserNode, generate.ParserHelm),
		// only provide Helm templates to generate only Helm part
		engine.WithTemplates(templates.FS(), templates.Helm()),
	),
}

func init() {
	rootCmd.AddCommand(chartCmd)
}

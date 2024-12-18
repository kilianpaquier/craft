package schemas

import "embed"

const (
	Chart = "chart.schema.json"
	Craft = "craft.schema.json"
)

//go:embed *.json
var fs embed.FS

// FS returns the embed.FS that contains schema files.
func FS() embed.FS {
	return fs
}

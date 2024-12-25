package schemas

import (
	"embed"
)

const (
	Chart = "chart.schema.json"
	Craft = "craft.schema.json"
)

//go:embed *.json
var fs embed.FS

func ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(name)
}

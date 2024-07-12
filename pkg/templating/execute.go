package templating

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	cfs "github.com/kilianpaquier/craft/pkg/fs"
)

// Execute runs tmpl.Execute with input data and write result into given dest file.
//
// When Execute is called, it deletes dest in case it already exists and reevaluate its rights (specific to linux).
func Execute(tmpl *template.Template, data any, dest string) error {
	// create destination directory only if one file would be generated
	if err := os.MkdirAll(filepath.Dir(dest), cfs.RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("create directory: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return fmt.Errorf("template execution: %w", err)
	}

	return cfs.WriteFile(dest, []byte(result.String()))
}

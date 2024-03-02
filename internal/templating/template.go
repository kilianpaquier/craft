package templating

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
)

// Execute runs Execute function from input tmpl with input data and write result to given dest file.
func Execute(tmpl *template.Template, dest string, data any) error {
	// create destination directory only if one file would be generated
	if err := os.MkdirAll(filepath.Dir(dest), filesystem.RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return fmt.Errorf("failed to template %s: %w", dest, err)
	}

	// check dest rights to apply (644 or 755)
	rights := func() fs.FileMode {
		if filepath.Ext(dest) == ".sh" {
			return filesystem.RwxRxRxRx
		}
		return filesystem.RwRR
	}()

	// remove file before rewritting it (in case rights changed)
	if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove %s before rewritting it: %w", dest, err)
	}

	// write new file content
	if err := os.WriteFile(dest, []byte(result.String()), rights); err != nil {
		return fmt.Errorf("failed to write %s: %w", dest, err)
	}
	return nil
}

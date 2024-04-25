package templating

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
)

// Execute runs Execute function from input tmpl with input data and write result to given dest file.
func Execute(tmpl *template.Template, data any, dest string) error {
	// create destination directory only if one file would be generated
	if err := os.MkdirAll(filepath.Dir(dest), filesystem.RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("create directory: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return fmt.Errorf("template execution: %w", err)
	}

	return WriteFile(dest, []byte(result.String()))
}

// WriteFile removes dest and rewrites it with input content.
// It's done as is to ensure file rights are recalculated.
func WriteFile(dest string, content []byte) error {
	// remove file before rewritting it (in case rights changed)
	if err := os.Remove(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("delete file: %w", err)
	}

	// write new file content
	if err := os.WriteFile(dest, content, GetRights(dest)); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// GetRights returns the appropriate file mode according to input file path.
func GetRights(dest string) os.FileMode {
	if strings.HasSuffix(dest, ".sh") {
		return filesystem.RwxRxRxRx
	}
	return filesystem.RwRR
}

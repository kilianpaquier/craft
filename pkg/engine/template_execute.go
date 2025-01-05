package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/kilianpaquier/craft/pkg/engine/files"
)

// Execute runs tmpl.Execute with input data and write result into given dest file.
//
// When Execute is called, it deletes dest in case it already exists and reevaluate its rights (specific to linux).
func execute(tmpl *template.Template, data any, dst string) error {
	// create destination directory only if one file would be generated
	if err := os.MkdirAll(filepath.Dir(dst), files.RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("create directory: %w", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return fmt.Errorf("template execution: %w", err)
	}

	if err := os.WriteFile(dst, result.Bytes(), files.RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	mode := files.RwRR
	if slices.Contains([]string{".sh"}, filepath.Ext(dst)) {
		mode = files.RwxRxRxRx
	}
	// force refresh rights since WriteFile doesn't do it
	// in case the target file already exists
	if err := os.Chmod(dst, mode); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}
	return nil
}

package templating

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
)

// Execute runs tmpl.Execute with input data and write result into given dest file.
//
// When Execute is called, it deletes dest in case it already exists and reevaluate its rights (specific to linux).
func Execute(tmpl *template.Template, data any, dest string) error {
	// create destination directory only if one file would be generated
	if err := os.MkdirAll(filepath.Dir(dest), cfs.RwxRxRxRx); err != nil && !os.IsExist(err) {
		return fmt.Errorf("create directory: %w", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return fmt.Errorf("template execution: %w", err)
	}

	if err := os.WriteFile(dest, result.Bytes(), cfs.RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	mode := cfs.RwRR
	if slices.Contains([]string{".sh"}, filepath.Ext(dest)) {
		mode = cfs.RwxRxRxRx
	}
	// force refresh rights since WriteFile doesn't do it
	// in case the target file already exists
	if err := os.Chmod(dest, mode); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}
	return nil
}

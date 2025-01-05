package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/files"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// ParserNode detects the presence of a ParserNode.js project by looking for a package.json file.
//
// In case of success, the function will set the language to "node"
// and the worker to "main" if the main property is present in the package.json file.
func ParserNode(ctx context.Context, destdir string, config *craft.Config) error {
	var jsonfile parser.PackageJSON
	jsonpath := filepath.Join(destdir, parser.FilePackageJSON)
	if err := files.ReadJSON(jsonpath, &jsonfile, os.ReadFile); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read json: %w", err)
	}
	engine.GetLogger(ctx).Infof("node detected, a '%s' is present and valid", parser.FilePackageJSON)

	if err := jsonfile.Validate(); err != nil {
		return fmt.Errorf("validate '%s': %w", parser.FilePackageJSON, err)
	}

	config.SetLanguage("node", jsonfile)
	if jsonfile.Main != nil {
		config.SetWorker("main")
	}
	return nil
}

var _ engine.Parser[craft.Config] = ParserNode // ensure interface is implemented

package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	craft "github.com/kilianpaquier/craft/pkg/craft/configuration"
	"github.com/kilianpaquier/craft/pkg/engine"
	"github.com/kilianpaquier/craft/pkg/engine/parser"
)

// ParserGolang detects the presence of a go.mod file
// and adds golang to languages.
//
// It also detects the presence of main.go files in cmd folder
// and adds them to executables composition.
//
// If a hugo config or theme file is present, it will be detected
// and hugo will be set as the language (golang will not in that cas).
func ParserGolang(ctx context.Context, destdir string, config *craft.Config) error {
	gomod, err := parser.ReadGomod(destdir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read '%s': %w", parser.FileGomod, err)
	}

	if hugoconfig, ok := parser.Hugo(destdir); ok {
		engine.GetLogger(ctx).Infof("hugo detected, theme or hugo files are present")
		config.SetLanguage("hugo", hugoconfig)
		return nil
	}

	engine.GetLogger(ctx).Infof("golang detected, file '%s' is present and valid", parser.FileGomod)
	config.SetLanguage("golang", gomod)

	executables, err := parser.ReadGoCmd(destdir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read '%s': %w", parser.FolderCMD, err)
	}

	config.Executables = executables
	return nil
}

var _ engine.Parser[craft.Config] = ParserGolang // ensure interface is implemented

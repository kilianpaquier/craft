package detectgen

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-playground/validator/v10"
	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/kilianpaquier/craft/internal/generate/filehandler"
	"github.com/kilianpaquier/craft/internal/models"
	"github.com/kilianpaquier/craft/internal/templating"
)

//go:generate go-builder-generator generate -f dir_generate.go -s DirGenerate -p set -d . --validate-func Validate

// generated is the regexp for generated files.
var generated = regexp.MustCompile(`Code generated [a-z-_0-9\ ]+; DO NOT EDIT\.`)

// GetGenerateFunc is a simplified function returning a basic GenerateFunc for an input ExecName.
//
// It uses behind the hood a private builder for dirGenerate which is the main function for all craft generations.
func GetGenerateFunc(name GenerateName) GenerateFunc {
	return func(ctx context.Context, config models.GenerateConfig, fsys filesystem.FS) error {
		generate, err := NewDirGenerateBuilder().
			SetConfig(config).
			SetData(config).
			SetFileHandlers(filehandler.AllHandlers(config)).
			SetFS(fsys).
			SetName(name).
			Build()
		if err != nil {
			return fmt.Errorf("invalid default generate build: %w", err)
		}
		return generate.Generate(ctx, config.Options.TemplatesDir, config.Options.DestinationDir)
	}
}

// DirGenerate represents the struct for basic dir generation.
//
// It takes various inputs that will never change during recursive calls
// to avoid passing them as input arguments.
type DirGenerate struct {
	Config       models.GenerateConfig `validate:"required"`
	Data         any                   `validate:"required"`
	FileHandlers []filehandler.Handler `validate:"omitempty,dive,required"`
	FS           filesystem.FS         `validate:"required"`
	Name         GenerateName          `validate:"required"`
}

// Validate ensures the built d is valid.
func (d *DirGenerate) Validate() error {
	if err := validator.New().Struct(d); err != nil {
		return fmt.Errorf("invalid dir generate build: %w", err)
	}
	return nil
}

// Generate walks over input srcdir and apply template of every src entry into destdir.
//
// If src entry is a directory and this directory name is the same as generate struct name then it dives into and executes Execute inside.
func (d *DirGenerate) Generate(ctx context.Context, srcdir, destdir string) error {
	// read source directory
	entries, err := d.FS.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	errs := lo.Map(entries, func(entry fs.DirEntry, _ int) error {
		src := path.Join(srcdir, entry.Name())
		if entry.IsDir() {
			// apply generation at root if the folder name is the dir generate name
			if entry.Name() == string(d.Name) {
				return d.Generate(ctx, src, destdir)
			}

			// apply templates on subdirs not being those associated to another generation
			if !slices.Contains(ReservedNames(), entry.Name()) {
				dest := filepath.Join(destdir, entry.Name())
				return d.Generate(ctx, src, dest)
			}
		}
		return d.handleFile(ctx, src, destdir, entry)
	})
	return errors.Join(errs...)
}

// handleFile is a private function used in Generate function to handle
// a specific file entry during iterative loops over folders and subfolders.
func (d *DirGenerate) handleFile(ctx context.Context, src, destdir string, entry fs.DirEntry) error {
	log := logrus.WithContext(ctx)

	// don't template files without .tmpl extension
	if !strings.HasSuffix(entry.Name(), models.TmplExtension) {
		return nil
	}
	filename := strings.TrimSuffix(entry.Name(), models.TmplExtension)
	dest := filepath.Join(destdir, filename)

	// verify that file matches generation rules
	generate := d.Config.Options.ForceAll || IsGenerated(dest) || slices.Contains(d.Config.Options.Force, filename)
	singleGeneration := filename == models.CraftFile && filesystem.Exists(dest)
	if !generate || singleGeneration {
		log.Warnf("not copying %s because it already exists", filename)
		return nil
	}

	// check if filename matches an optional file
	for _, handler := range d.FileHandlers {
		ok, apply := handler(src, dest, filename)
		if !ok {
			continue
		}
		if apply {
			break // break loop since optional file was found
		}

		// handle optional file deletion
		if err := os.RemoveAll(dest); err != nil && !errors.Is(err, fs.ErrNotExist) {
			log.WithError(err).Warn("failed to delete non applicable file")
		}
		return nil
	}

	tmpl, err := template.New(entry.Name()).
		Funcs(sprig.FuncMap()).
		Funcs(templating.FuncMap()).
		Delims(d.Config.Options.StartDelim, d.Config.Options.EndDelim).
		ParseFS(d.FS, src)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", src, err)
	}
	return templating.Execute(tmpl, d.Data, dest)
}

// IsGenerated returns truthy if input destination is a generated file.
func IsGenerated(dest string) bool {
	// retrieve file content, if there's an error, generation to make
	content, err := os.ReadFile(dest)
	if err != nil {
		return true
	}

	// never regenerate craft files (it's handled differently)
	if filepath.Base(dest) == models.CraftFile {
		return false
	}

	// special case (shouldn't happen) where the destination has been replaced with an empty file
	if len(content) == 0 {
		return true
	}
	lines := strings.Split(string(content), "\n")

	// check first line for generated regexp
	if len(lines) >= 1 && generated.Match([]byte(lines[0])) {
		return true
	}

	// check second line for generated regexp
	if len(lines) >= 2 && generated.Match([]byte(lines[1])) {
		return true
	}
	return false
}

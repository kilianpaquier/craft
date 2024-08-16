package generate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"github.com/kilianpaquier/cli-sdk/pkg/clog"

	"github.com/kilianpaquier/craft/pkg/craft"
)

// DetectHelm handles the detection of helm chart generation option in metadata
// and returns the appropriate slice of Exec.
func DetectHelm(_ context.Context, log clog.Logger, _ string, metadata Metadata) (Metadata, []Exec, error) {
	if metadata.NoChart {
		return metadata, []Exec{removeHelm}, nil
	}
	log.Infof("helm chart detected, %s doesn't have no_chart key", craft.File)
	return metadata, []Exec{generateHelm}, nil
}

var _ Detect = DetectHelm // ensure interface is implemented

// generateHelm generates the appropriate helm chart at destdir.
//
// To be able to use the maximum number of variables in templates (in input fsys inside helm folder),
// a marshal is applied on input config and on {{destdir}}/chart/.craft.
func generateHelm(_ context.Context, log clog.Logger, fsys cfs.FS, srcdir, destdir string, metadata Metadata, opts ExecOpts) error { // nolint:revive
	srcdir = path.Join(srcdir, "lang_helm")   // nolint:revive
	destdir = filepath.Join(destdir, "chart") // nolint:revive

	// transform craft configuration into generic chart configuration (easier to maintain)
	var chart map[string]any
	bytes, _ := json.Marshal(metadata)
	_ = json.Unmarshal(bytes, &chart)

	// read overrides values
	var overrides map[string]any
	if err := craft.Read(destdir, &overrides); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("read helm chart overrides: %w", err)
	}

	// merge overrides into chart with overwrite
	if err := mergo.Merge(&chart, overrides, mergo.WithOverride); err != nil {
		return fmt.Errorf("merge helm chart overrides with craft configuration: %w", err)
	}

	return handleDir(log, fsys, srcdir, destdir, chart, "helm", opts)
}

// removeHelm deletes the chart folder inside destdir.
func removeHelm(_ context.Context, _ clog.Logger, _ cfs.FS, _, destdir string, _ Metadata, _ ExecOpts) error { // nolint:revive
	if err := os.RemoveAll(filepath.Join(destdir, "chart")); err != nil {
		return fmt.Errorf("delete directory: %w", err)
	}
	return nil
}

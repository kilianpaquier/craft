package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"gopkg.in/yaml.v3"

	"github.com/kilianpaquier/craft/internal/models"
)

// ReadCraft reads the .craft file in srcdir input into the out input.
func ReadCraft(srcdir string, out any) error {
	craft := filepath.Join(srcdir, models.CraftFile)

	content, err := os.ReadFile(craft)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fs.ErrNotExist
		}
		return fmt.Errorf("read file: %w", err)
	}

	if err := yaml.Unmarshal(content, out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

// WriteCraft writes the input craft into the input destdir in .craft file.
func WriteCraft(destdir string, config models.CraftConfig) error {
	craft := filepath.Join(destdir, models.CraftFile)

	// create a buffer with craft notice
	buffer := bytes.NewBuffer([]byte("# Craft configuration file (https://github.com/kilianpaquier/craft)\n---\n"))

	// create yaml encoder and writes the full configuration in the buffer,
	// following the craft notice
	encoder := yaml.NewEncoder(buffer)
	defer encoder.Close()
	encoder.SetIndent(2)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("encode file: %w", err)
	}

	if err := os.WriteFile(craft, buffer.Bytes(), filesystem.RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

// EnsureDefaults acts to ensure default properties are always sets and also migrates old properties into new fields.
func EnsureDefaults(config models.CraftConfig) models.CraftConfig {
	if config.CI != nil {
		// sets default release mode for github actions
		if config.CI.Name == models.Github && config.CI.Release.Mode == "" {
			config.CI.Release.Mode = models.GithubToken
		}

		// keep release mode empty when working with gitlab CICD
		if config.CI.Name == models.Gitlab && config.CI.Release.Mode != "" {
			config.CI.Release.Mode = ""
		}

		// migrate old auto_release option
		if slices.Contains(config.CI.Options, "auto_release") {
			config.CI.Release.Auto = true
			config.CI.Options = slices.DeleteFunc(config.CI.Options, func(option string) bool { return option == "auto_release" })
		}

		// migrate old backmerge optin
		if slices.Contains(config.CI.Options, "backmerge") {
			config.CI.Release.Backmerge = true
			config.CI.Options = slices.DeleteFunc(config.CI.Options, func(option string) bool { return option == "backmerge" })
		}
	}
	return config
}

package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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

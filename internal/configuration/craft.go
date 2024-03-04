package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	filesystem "github.com/kilianpaquier/filesystem/pkg"
	"gopkg.in/yaml.v3"

	"github.com/kilianpaquier/craft/internal/models"
)

// ReadCraft reads the .craft file in srcdir input into the out input.
func ReadCraft(srcdir string, out any) error {
	craft := filepath.Join(srcdir, models.CraftFile)
	bytes, err := os.ReadFile(craft)
	if err != nil {
		if os.IsNotExist(err) {
			return os.ErrNotExist
		}
		return fmt.Errorf("failed to read file: %w", err)
	}
	if err := yaml.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}
	return nil
}

// WriteCraft writes the input craft into the input destdir in .craft file.
func WriteCraft(destdir string, config models.CraftConfig) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	craft := filepath.Join(destdir, models.CraftFile)
	return os.WriteFile(craft, bytes, filesystem.RwRR)
}

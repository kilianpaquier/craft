package craft

import (
	"bytes"
	"fmt"
	"os"

	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
	"gopkg.in/yaml.v3"
)

// Read reads the .craft file in srcdir input into the out input.
func Read(src string, out any) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if err := yaml.Unmarshal(content, out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

// Write writes the input craft into the input destdir in .craft file.
func Write(dest string, config Configuration) error {
	// create a buffer with craft notice
	buffer := bytes.NewBufferString("# Craft configuration file (https://github.com/kilianpaquier/craft)\n---\n")

	// create yaml encoder and writes the full configuration in the buffer,
	// following the craft notice
	encoder := yaml.NewEncoder(buffer)
	defer encoder.Close()
	encoder.SetIndent(2)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("encode file: %w", err)
	}

	if err := os.WriteFile(dest, buffer.Bytes(), cfs.RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

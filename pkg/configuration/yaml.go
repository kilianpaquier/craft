package configuration

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/kilianpaquier/cli-sdk/pkg/cfs"
)

// ReadYAML reads the input src and unmarshal it into the out configuration.
func ReadYAML(src string, out any) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	if err := yaml.Unmarshal(content, out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

// WriteYAML writes the input configuration into the dest.
func WriteYAML(dest string, config any, opts ...yaml.EncodeOption) error {
	bytes, err := yaml.MarshalWithOptions(config, opts...)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(dest, bytes, cfs.RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

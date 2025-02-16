package files

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// ReadYAML reads the input src
// and unmarshal it in YAML into the out configuration.
func ReadYAML(src string, out any, read func(src string) ([]byte, error)) error {
	if read == nil {
		return ErrNilRead
	}

	content, err := read(src)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	if err := yaml.Unmarshal(content, out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

// WriteYAML writes the input configuration into the dest in YAML format.
func WriteYAML(out string, data any, opts ...yaml.EncodeOption) error {
	content, err := yaml.MarshalWithOptions(data, opts...)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(out, content, RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

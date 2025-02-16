package files

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSON reads the input src
// and unmarshal it in JSON format into the out configuration.
func ReadJSON(src string, out any, read func(src string) ([]byte, error)) error {
	if read == nil {
		return ErrNilRead
	}

	content, err := read(src)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	if err := json.Unmarshal(content, out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

// WriteJSON writes the input configuration into the dest in JSON format.
func WriteJSON(out string, data any) error {
	content, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(out, content, RwRR); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

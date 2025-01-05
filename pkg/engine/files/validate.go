package files

import (
	"errors"
	"fmt"
	"path"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var defaultPrinter = message.NewPrinter(language.English)

// ErrNilRead is returned when the read function is nil.
var ErrNilRead = errors.New("read function is nil")

// Validate validates a given file (read from readFile) with a JSON schema (read from readSchema).
//
// This function takes both delegated functions to be able to read from various fs.
//
// Example:
//
//	func main() {
//		err := files.Validate(
//			func(out any) error { return files.ReadJSON("path/to/schema", out, os.ReadFile) },
//			func(out any) error { return files.ReadYAML("path/to/file/to/validate", out, os.ReadFile) },
//		)
//		// handle err
//	}
func Validate(readSchema, readFile func(v any) error) error {
	if readSchema == nil || readFile == nil {
		return ErrNilRead
	}

	var schema any
	if err := readSchema(&schema); err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	_ = compiler.AddResource("schema.json", schema)

	sch, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	var doc any
	if err := readFile(&doc); err != nil {
		return err // error is already wrapped
	}

	if err := sch.Validate(doc); err != nil {
		ve := &jsonschema.ValidationError{}
		if errors.As(err, &ve) {
			return fmt.Errorf("validate schema:\n%w", errors.Join(flatten(ve)...))
		}
		return fmt.Errorf("validate schema: %w", err)
	}
	return nil
}

// ValidationError represents a simplified view of jsonschema.ValidationError.
//
// It it used to override specific error messages (like kind.FalseSchema "false schema") in craft validation context.
type ValidationError struct {
	Message  string
	Property string
}

var _ error = &ValidationError{}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("- at '%s': %s", v.Property, v.Message)
}

// flatten converts a jsonschema.ValidationError to a ValidationError.
func flatten(ve *jsonschema.ValidationError) []error {
	var errs []error
	if len(ve.Causes) == 0 {
		property := "/" + path.Join(ve.InstanceLocation...)
		switch ek := ve.ErrorKind.(type) {
		case *kind.FalseSchema:
			errs = append(errs, &ValidationError{
				Property: property,
				Message:  "must not be provided",
			})
		default:
			errs = append(errs, &ValidationError{
				Property: property,
				Message:  ek.LocalizedString(defaultPrinter),
			})
		}
	}

	for _, cause := range ve.Causes {
		errs = append(errs, flatten(cause)...)
	}
	return errs
}

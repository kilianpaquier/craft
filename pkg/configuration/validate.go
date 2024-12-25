package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var defaultPrinter = message.NewPrinter(language.English)

// Validate validates input src following the input schema (read during execution).
func Validate(src string, readSchema func() ([]byte, error)) error {
	bytes, err := readSchema()
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	var schema any
	if err := json.Unmarshal(bytes, &schema); err != nil {
		return fmt.Errorf("unmarshal schema: %w", err)
	}
	compiler := jsonschema.NewCompiler()
	_ = compiler.AddResource("schema.json", schema)

	sch, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	var doc any
	if err := ReadYAML(src, &doc); err != nil {
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
		switch ve.ErrorKind.(type) {
		// case *kind.AllOf, *kind.AnyOf, *kind.OneOf, *kind.Group, *kind.Schema:
		// 	err = &ValidationError{}
		case *kind.FalseSchema:
			errs = append(errs, &ValidationError{
				Property: property,
				Message:  "must not be provided",
			})
		default:
			errs = append(errs, &ValidationError{
				Property: property,
				Message:  ve.ErrorKind.LocalizedString(defaultPrinter),
			})
		}
	}

	for _, cause := range ve.Causes {
		errs = append(errs, flatten(cause)...)
	}
	return errs
}

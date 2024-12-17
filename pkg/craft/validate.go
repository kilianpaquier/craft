package craft

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	schemas "github.com/kilianpaquier/craft/.schemas"
)

var defaultPrinter = message.NewPrinter(language.English)

// Validate validates .craft file from srcdir following craft schema.
func Validate(src string) error {
	bytes, err := schemas.FS().ReadFile(schemas.Craft)
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	var schema any
	if err := json.Unmarshal(bytes, &schema); err != nil {
		return fmt.Errorf("unmarshal schema: %w", err)
	}
	compiler := jsonschema.NewCompiler()
	_ = compiler.AddResource(schemas.Craft, schema)

	sch, err := compiler.Compile(schemas.Craft)
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	var doc any
	if err := Read(src, &doc); err != nil {
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

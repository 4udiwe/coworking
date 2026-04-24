package json_schema_validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Validator struct {
	schema *jsonschema.Schema
}

func NewValidator(schemaData string) (*Validator, error) {
	compiler := jsonschema.NewCompiler()

	if err := compiler.AddResource("schema.json", strings.NewReader(schemaData)); err != nil {
		return nil, err
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, err
	}

	return &Validator{
		schema: schema,
	}, nil
}

func (v *Validator) Validate(jsonData []byte) error {
	var data interface{}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := v.schema.Validate(data); err != nil {
		if ve, ok := err.(*jsonschema.ValidationError); ok {
			return fmt.Errorf("schema validation failed: %s", ve.Error())
		}
		return err
	}

	return nil
}

package app

import (
	layout_schema "github.com/4udiwe/cowoking/booking-service/internal/layout/schema"
	"github.com/4udiwe/cowoking/booking-service/pkg/json_schema_validator"
)

func (app *App) LayoutValidator() *json_schema_validator.Validator {
	if app.layoutValidator != nil {
		return app.layoutValidator
	}
	app.layoutValidator, _ = json_schema_validator.NewValidator(layout_schema.LayoutSchemaData)
	return app.layoutValidator
}

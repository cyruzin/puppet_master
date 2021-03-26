package validation

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// APIMessage is a struct for generic JSON response.
type APIMessage struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (a *APIMessage) Error() string {
	return a.Message
}

func validationMap(err validator.FieldError, field string) *APIMessage {
	errMap := map[string]string{
		"required": "is required",
		"email":    "is not valid",
		"min":      "minimum length is " + err.Param(),
		"gte":      "minimum length is " + err.Param(),
	}

	currentField := strings.ToLower(err.Field())

	fmt.Println(field)

	if field != "" {
		currentField = field
	}

	return &APIMessage{
		Message: "the " + currentField + " field " + errMap[err.Tag()],
	}
}

// validatorMessage handles validation error messages.
func validatorMessage(err error, field string) error {
	apiErrors := err.(validator.ValidationErrors)

	return validationMap(apiErrors[0], field)
}

// IsAValidSchema checks if a given schema is valid.
func IsAValidSchema(ctx context.Context, schema interface{}) error {
	validate := validator.New()

	if err := validate.StructCtx(ctx, schema); err != nil {
		return validatorMessage(err, "")
	}

	return nil
}

// IsAValidField checks if a given schema is field.
func IsAValidField(ctx context.Context, field interface{}, name, tag string) error {
	validate := validator.New()

	if err := validate.VarCtx(ctx, field, tag); err != nil {
		return validatorMessage(err, name)
	}

	return nil
}

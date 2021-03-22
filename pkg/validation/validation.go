package validation

import (
	"context"
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

func validationMap(err validator.FieldError) *APIMessage {
	errMap := map[string]string{
		"required": "is required",
		"email":    "is not valid",
		"min":      "minimum length is " + err.Param(),
		"gte":      "minimum length is " + err.Param(),
	}

	return &APIMessage{
		Message: "The " + strings.ToLower(err.Field()) + " field " + errMap[err.Tag()],
	}
}

// validatorMessage handles validation error messages.
func validatorMessage(err error) error {
	apiErrors := err.(validator.ValidationErrors)

	return validationMap(apiErrors[0])
}

// IsAValidSchema checks if a given schema is valid.
func IsAValidSchema(ctx context.Context, schema interface{}) error {
	validate := validator.New()

	if err := validate.StructCtx(ctx, schema); err != nil {
		return validatorMessage(err)
	}

	return nil
}

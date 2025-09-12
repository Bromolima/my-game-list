package validation

import (
	"encoding/json"
	"errors"

	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (c *CustomValidator) Validate(i any) error {
	if err := c.validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	return nil
}

func ValidateUserError(validationErr error) *resterr.RestErr {
	var jsonErr *json.UnmarshalTypeError
	var jsonValidationError validator.ValidationErrors

	if errors.As(validationErr, &jsonErr) {
		return resterr.NewBadRequestError("Invalid field type")
	}

	if errors.As(validationErr, &jsonValidationError) {
		errorsCauses := []resterr.Causes{}

		for _, e := range validationErr.(validator.ValidationErrors) {
			cause := resterr.Causes{
				Message: e.Translate(Transl),
				Field:   e.Field(),
			}

			errorsCauses = append(errorsCauses, cause)
		}

		return resterr.NewBadRequestValidationError("Some fields are invalid", errorsCauses)
	}

	return resterr.NewBadRequestError("Error trying to convert field")
}

package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator struct {
	validator *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func New() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	registerCustomValidators(v)

	return &Validator{
		validator: v,
	}
}

func (v *Validator) Validate(s interface{}) map[string]string {
	errors := make(map[string]string)

	if err := v.validator.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			errors[fieldName] = v.getErrorMessage(err)
		}
	}

	return errors
}

func (v *Validator) ParseAndValidate(c *fiber.Ctx, dest interface{}) error {
	if err := c.BodyParser(dest); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	if errors := v.Validate(dest); len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}

	return nil
}

type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}

func (ve *ValidationErrors) Error() string {
	return "validation failed"
}

func (v *Validator) getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters long", field, param)
		}
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must not exceed %s characters", field, param)
		}
		return fmt.Sprintf("%s must not exceed %s", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must be a valid number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "password":
		return fmt.Sprintf("%s must contain at least 8 characters with uppercase, lowercase, number and special character", field)
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "username":
		return fmt.Sprintf("%s must be 3-30 characters and contain only letters, numbers, dots, underscores and hyphens", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func registerCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}

		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

		return hasUpper && hasLower && hasNumber && hasSpecial
	})

	//custom validators
}

package middleware

import (
	"cms-octo-chat-api/helper"
	bmodel "cms-octo-chat-api/model"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

type ValidationMiddleware struct {
	validate   *validator.Validate
	msgMapping map[string]string
}

// Load YAML error messages
func loadValidationMessages(filePath string) (map[string]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var messages map[string]string
	if err := yaml.Unmarshal(data, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// Create new middleware instance
func (m *BaseMiddleware) NewBodyValidatorMiddleware(yamlPath string) (*ValidationMiddleware, error) {
	messages, err := loadValidationMessages(yamlPath)
	if err != nil {
		return nil, err
	}
	// initiate validator
	newValidator := validator.New(validator.WithRequiredStructEnabled())

	// register custom validator here
	newValidator.RegisterValidation("exists", helper.ExistsValidator(m.DB))

	return &ValidationMiddleware{
		validate:   newValidator,
		msgMapping: messages,
	}, nil
}

// Middleware function: Accepts a pointer to a struct and uses reflection to create new instances per request
func (vm *ValidationMiddleware) Validate(model interface{}) fiber.Handler {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return func(c *fiber.Ctx) error {
		// Create a new instance of the model for each request
		dest := reflect.New(modelType).Interface()

		// Parse body into struct
		if err := c.BodyParser(dest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(bmodel.ErrorValidationResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		// Validate struct
		if err := vm.validate.Struct(dest); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				customErrors := make(map[string]string)
				for _, e := range errs {
					tag := e.Tag()
					field := e.Field()
					param := e.Param()

					// Default message if YAML doesn't have it
					msgTemplate := vm.msgMapping[tag]
					if msgTemplate == "" {
						msgTemplate = fmt.Sprintf("%s failed on the '%s' rule", field, tag)
					}

					// Replace placeholders
					msg := strings.ReplaceAll(msgTemplate, "{field}", field)
					msg = strings.ReplaceAll(msg, "{param}", param)

					customErrors[strings.ToLower(field)] = msg
				}

				return c.Status(fiber.StatusBadRequest).JSON(bmodel.ErrorValidationResponse{
					Code:    fiber.StatusBadRequest,
					Message: "Validation failed",
					Errors:  customErrors,
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(bmodel.ErrorValidationResponse{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}

		// Save validated struct for handler
		c.Locals("validatedBody", dest)
		return c.Next()
	}
}

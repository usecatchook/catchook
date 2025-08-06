package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      *ErrorData  `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

type ErrorData struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Field   string            `json:"field,omitempty"`
}

type Pagination struct {
	CurrentPage int  `json:"current_page,omitempty"`
	TotalPages  int  `json:"total_pages,omitempty"`
	Total       int  `json:"total,omitempty"`
	Limit       int  `json:"limit,omitempty"`
	HasNext     bool `json:"has_next,omitempty"`
	HasPrev     bool `json:"has_prev,omitempty"`
}

type ValidationError struct {
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	Errors    map[string]string `json:"errors"`
	Timestamp time.Time         `json:"timestamp"`
}

func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func Paginated(c *fiber.Ctx, data interface{}, pagination Pagination, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: &pagination,
		Timestamp:  time.Now().UTC(),
	})
}

func BadRequest(c *fiber.Ctx, message string, details map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "BAD_REQUEST",
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
	})
}

func ValidationFailed(c *fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(ValidationError{
		Success:   false,
		Message:   "Validation failed",
		Errors:    errors,
		Timestamp: time.Now().UTC(),
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "FORBIDDEN",
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "NOT_FOUND",
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}

func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "CONFLICT",
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}

func UnprocessableEntity(c *fiber.Ctx, message string, field string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "UNPROCESSABLE_ENTITY",
			Message: message,
			Field:   field,
		},
		Timestamp: time.Now().UTC(),
	})
}

func TooManyRequests(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "TOO_MANY_REQUESTS",
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}

func InternalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "INTERNAL_ERROR",
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}

func ServiceUnavailable(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusServiceUnavailable).JSON(Response{
		Success: false,
		Error: &ErrorData{
			Code:    "SERVICE_UNAVAILABLE",
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
}

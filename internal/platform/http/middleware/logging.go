package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/theotruvelot/catchook/pkg/logger"
)

func RequestLogging(baseLogger logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := generateRequestID()

		// Store request_id in Fiber context
		fiberCtx := context.WithValue(c.UserContext(), logger.RequestIDKey, requestID)
		c.SetUserContext(fiberCtx)

		fields := []zap.Field{
			logger.String("method", c.Method()),
			logger.String("path", c.Path()),
			logger.String("ip", c.IP()),
		}

		if c.Method() != "GET" && c.Method() != "DELETE" {
			if body := c.Body(); len(body) > 0 && len(body) < 1024 {
				fields = append(fields, logger.String("body", maskSensitiveData(string(body))))
			}
		}

		baseLogger.Info(fiberCtx, "Request received", fields...)

		err := c.Next()

		duration := time.Since(start)
		responseFields := []zap.Field{
			logger.String("method", c.Method()),
			logger.String("path", c.Path()),
			logger.Int("status", c.Response().StatusCode()),
			logger.Duration("duration", duration.Milliseconds()),
			logger.Int("response_size", len(c.Response().Body())),
		}

		if userID, err := GetAuthUserID(c); err == nil {
			responseFields = append(responseFields, logger.String("user_id", userID))
		}

		if err != nil {
			responseFields = append(responseFields, logger.Error(err))
			baseLogger.Error(fiberCtx, "Request failed", responseFields...)
		} else if c.Response().StatusCode() >= 500 {
			baseLogger.Error(fiberCtx, "Request completed with server error", responseFields...)
		} else if c.Response().StatusCode() >= 400 {
			baseLogger.Warn(fiberCtx, "Request completed with client error", responseFields...)
		} else {
			baseLogger.Info(fiberCtx, "Request completed", responseFields...)
		}

		return err
	}
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetContextWithRequestID returns a Go context with the request_id from Fiber context
func GetContextWithRequestID(c *fiber.Ctx) context.Context {
	base := c.UserContext()
	if requestID := base.Value(logger.RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return context.WithValue(base, logger.RequestIDKey, id)
		}
	}
	return base
}

func maskSensitiveData(body string) string {
	sensitiveFields := []string{"password", "token", "secret", "authorization", "apikey", "api_key"}
	masked := body

	for _, field := range sensitiveFields {
		pattern := `("` + field + `")\s*:\s*(".*?")`
		re := regexp.MustCompile(`(?i)` + pattern) // (?i) pour ignorer la casse

		masked = re.ReplaceAllString(masked, `$1: "***"`)
	}

	return masked
}

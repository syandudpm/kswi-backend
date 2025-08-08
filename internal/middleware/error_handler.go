package middleware

import (
	"kswi-backend/internal/config"
	"kswi-backend/internal/shared/errors"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(cfg *config.Config) gin.HandlerFunc {
	isProduction := cfg.App.Environment == "production"

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last().Err
			var appErr *errors.AppError

			// Convert to AppError if not already one
			if !errors.As(lastErr, &appErr) {
				appErr = &errors.AppError{
					Type:    errors.TypeInternal,
					Message: "Internal server error",
				}

				if !isProduction {
					appErr.Details = lastErr.Error()
				}
			}

			// Enhanced logging using structured logger
			logger := config.GetSugaredLogger()

			// Base fields for all logs
			baseFields := []interface{}{
				"error_type", appErr.Type,
				"error_message", appErr.Message,
				"request_method", c.Request.Method,
				"request_path", c.Request.URL.Path,
				"request_ip", c.ClientIP(),
				"status_code", appErr.GetStatusCode(),
				"timestamp", time.Now().Format(time.RFC3339),
			}

			if !isProduction {
				// Development: Log everything with full details for debugging
				allFields := append(baseFields,
					"error_details", appErr.Details,
					"original_error", func() string {
						if appErr.GetOriginalError() != nil {
							return appErr.GetOriginalError().Error()
						}
						return lastErr.Error()
					}(),
					"request_headers", c.Request.Header,
					"user_agent", c.Request.UserAgent(),
					"gin_errors", c.Errors.String(),
				)

				// Add query params if present
				if len(c.Request.URL.RawQuery) > 0 {
					allFields = append(allFields, "query_params", c.Request.URL.RawQuery)
				}

				logger.Errorw("Request error occurred", allFields...)

			} else {
				// Production: Log with essential details but no sensitive info
				essentialFields := append(baseFields,
					"error_details", appErr.Details,
					"original_error", func() string {
						if appErr.GetOriginalError() != nil {
							return appErr.GetOriginalError().Error()
						}
						return lastErr.Error()
					}(),
					"user_agent", c.Request.UserAgent(),
				)

				// Only add request ID if present (for tracing)
				if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
					essentialFields = append(essentialFields, "request_id", requestID)
				}

				logger.Errorw("Request error occurred", essentialFields...)
			}

			// Prepare safe response
			response := errors.Response{
				Success: false,
				Error:   appErr,
			}

			// Sanitize for production
			if isProduction && errors.ShouldHideDetails(appErr.Type) {
				response.Error = &errors.AppError{
					Type:    errors.TypeInternal,
					Message: "Internal server error",
				}
			}

			c.AbortWithStatusJSON(appErr.GetStatusCode(), response)
		}
	}
}

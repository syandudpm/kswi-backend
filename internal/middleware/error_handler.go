package middleware

import (
	"kswi-backend/internal/config"
	"kswi-backend/internal/shared/errors"
	"time"

	"github.com/gin-gonic/gin"
)

// getLogLevel determines the appropriate log level based on error type and environment
func getLogLevel(appErr *errors.AppError, isProduction bool) string {
	switch appErr.Type {
	case errors.TypeInternal:
		return "error"

	case errors.TypeDatabase:
		return "error"

	case errors.TypeAuth:
		if isProduction {
			return "warn" // Don't flood error logs with failed login attempts
		}
		return "error"

	case errors.TypeNotFound:
		if isProduction {
			return "info" // 404s are usually not critical in production
		}
		return "warn"

	case errors.TypeValidation:
		if isProduction {
			return "info" // User input errors are not system errors
		}
		return "warn"

	default:
		return "error"
	}
}

// shouldSkipLogging determines if error should be skipped based on thresholds
func shouldSkipLogging(appErr *errors.AppError, isProduction bool, requestPath string) bool {
	if !isProduction {
		return false // Log everything in development
	}

	// In production, you might want to skip certain errors to reduce noise
	switch appErr.Type {
	case errors.TypeNotFound:
		// Skip logging 404s for common bot requests
		commonBotPaths := []string{"/favicon.ico", "/robots.txt", "/.env", "/wp-admin", "/admin"}
		for _, path := range commonBotPaths {
			if requestPath == path {
				return true
			}
		}
		return false

	case errors.TypeValidation:
		// Could add logic here to skip validation errors from suspicious IPs
		return false

	default:
		return false
	}
}

// getSafeHeaders returns only safe headers for production logging
func getSafeHeaders(c *gin.Context) map[string]string {
	safeHeaders := map[string]string{}

	safeHeaderKeys := []string{
		"Content-Type",
		"Content-Length",
		"Accept",
		"User-Agent",
		"X-Request-ID",
		"X-Trace-ID",
		"X-Correlation-ID",
		"Referer",
		"Origin",
	}

	for _, key := range safeHeaderKeys {
		if value := c.GetHeader(key); value != "" {
			safeHeaders[key] = value
		}
	}

	return safeHeaders
}

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

			// Check if we should skip logging this error
			if shouldSkipLogging(appErr, isProduction, c.Request.URL.Path) {
				c.AbortWithStatusJSON(appErr.GetStatusCode(), errors.Response{
					Success: false,
					Error:   appErr,
				})
				return
			}

			// Get your existing logger
			logger := config.GetSugaredLogger()

			// Determine log level
			logLevel := getLogLevel(appErr, isProduction)

			// Base fields for all logs
			baseFields := []interface{}{
				"error_type", appErr.Type,
				"error_message", appErr.Message,
				"log_level", logLevel,
				"request_method", c.Request.Method,
				"request_path", c.Request.URL.Path,
				"request_ip", c.ClientIP(),
				"status_code", appErr.GetStatusCode(),
				"timestamp", time.Now().Format(time.RFC3339),
			}

			if !isProduction {
				// Development: Log with full details
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

				if len(c.Request.URL.RawQuery) > 0 {
					allFields = append(allFields, "query_params", c.Request.URL.RawQuery)
				}

				// Use appropriate log level
				switch logLevel {
				case "error":
					logger.Errorw("Request error occurred", allFields...)
				case "warn":
					logger.Warnw("Request warning occurred", allFields...)
				case "info":
					logger.Infow("Request info occurred", allFields...)
				default:
					logger.Errorw("Request error occurred", allFields...)
				}

			} else {
				// Production: Log based on error severity
				essentialFields := append(baseFields,
					"user_agent", c.Request.UserAgent(),
				)

				// Add detailed info only for ERROR level logs
				if logLevel == "error" {
					essentialFields = append(essentialFields,
						"error_details", appErr.Details,
						"original_error", func() string {
							if appErr.GetOriginalError() != nil {
								return appErr.GetOriginalError().Error()
							}
							return lastErr.Error()
						}(),
						"safe_headers", getSafeHeaders(c),
					)
				}

				// Add request ID if present
				if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
					essentialFields = append(essentialFields, "request_id", requestID)
				}

				// Use appropriate log level
				switch logLevel {
				case "error":
					logger.Errorw("Request error occurred", essentialFields...)
				case "warn":
					logger.Warnw("Request warning occurred", essentialFields...)
				case "info":
					logger.Infow("Request info occurred", essentialFields...)
				default:
					logger.Errorw("Request error occurred", essentialFields...)
				}
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

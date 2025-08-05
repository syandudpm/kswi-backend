package middleware

import (
	"kswi-backend/internal/config"
	"kswi-backend/internal/shared/errors"

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
			if !isProduction {
				logger.Errorf("Request error: %+v", lastErr)
			} else if appErr.Type == errors.TypeInternal {
				logger.Errorf("Internal error: %v", lastErr)
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

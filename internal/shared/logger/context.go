package logger

import (
	"context"
	"kswi-backend/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ContextualLogger wraps zap logger with contextual information
type ContextualLogger struct {
	logger *zap.SugaredLogger
	fields []interface{}
}

// NewContextualLogger creates a logger with base context
func NewContextualLogger() *ContextualLogger {
	return &ContextualLogger{
		logger: config.GetSugaredLogger(),
		fields: make([]interface{}, 0),
	}
}

// WithRequestID adds request ID to logger context
func (cl *ContextualLogger) WithRequestID(requestID string) *ContextualLogger {
	if requestID == "" {
		return cl
	}

	newFields := append(cl.fields, "request_id", requestID)
	return &ContextualLogger{
		logger: cl.logger,
		fields: newFields,
	}
}

// WithModule adds module name to logger context
func (cl *ContextualLogger) WithModule(module string) *ContextualLogger {
	newFields := append(cl.fields, "module", module)
	return &ContextualLogger{
		logger: cl.logger,
		fields: newFields,
	}
}

// WithUserContext adds user-related context
func (cl *ContextualLogger) WithUserContext(userID, username string) *ContextualLogger {
	newFields := cl.fields
	if userID != "" {
		newFields = append(newFields, "user_id", userID)
	}
	if username != "" {
		newFields = append(newFields, "username", username)
	}

	return &ContextualLogger{
		logger: cl.logger,
		fields: newFields,
	}
}

// WithFields adds custom fields
func (cl *ContextualLogger) WithFields(keyValuePairs ...interface{}) *ContextualLogger {
	newFields := append(cl.fields, keyValuePairs...)
	return &ContextualLogger{
		logger: cl.logger,
		fields: newFields,
	}
}

// WithError adds error context for consistent error logging
func (cl *ContextualLogger) WithError(err error) *ContextualLogger {
	if err == nil {
		return cl
	}

	newFields := append(cl.fields, "error", err.Error())
	return &ContextualLogger{
		logger: cl.logger,
		fields: newFields,
	}
}

// Logging methods with performance optimization
func (cl *ContextualLogger) Debug(msg string) {
	if cl.logger.Desugar().Core().Enabled(zap.DebugLevel) {
		cl.logger.Debugw(msg, cl.fields...)
	}
}

func (cl *ContextualLogger) Info(msg string) {
	cl.logger.Infow(msg, cl.fields...)
}

func (cl *ContextualLogger) Warn(msg string) {
	cl.logger.Warnw(msg, cl.fields...)
}

func (cl *ContextualLogger) Error(msg string) {
	cl.logger.Errorw(msg, cl.fields...)
}

// FromGinContext creates a contextual logger from gin context
func FromGinContext(c *gin.Context) *ContextualLogger {
	logger := NewContextualLogger()

	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			logger = logger.WithRequestID(id)
		}
	}

	return logger
}

// FromContext creates a contextual logger from standard context
func FromContext(ctx context.Context) *ContextualLogger {
	logger := NewContextualLogger()

	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			logger = logger.WithRequestID(id)
		}
	}

	return logger
}

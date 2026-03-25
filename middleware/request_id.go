package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "requestID"

// generates or extracts a unique request ID for tracing
func RequestID(context *gin.Context) {
	requestID := context.GetHeader("X-Request-ID")

	if requestID == "" {
		requestID = uuid.New().String()
	}
	//store requestID in context
	context.Set(RequestIDKey, requestID)

	context.Writer.Header().Set("X-Request-ID", requestID)

	context.Next()
}

// retrieves the request ID from context
func GetRequestID(context *gin.Context) string {
	if requestID, exists := context.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return "unknown"
}

package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()
		method := context.Request.Method
		path := context.Request.URL.Path

		context.Next()

		status := context.Writer.Status()
		duration := time.Since(start)
		clientIP := context.ClientIP()

		// color code the status
		statusColor := colorForStatus(status)
		methodColor := colorForMethod(method)
		reset := "\033[0m"

		fmt.Printf("[API] %v | %s%d%s | %v | %s | %s%s%s %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusColor, status, reset,
			duration,
			clientIP,
			methodColor, method, reset,
			path,
		)
	}
}

func colorForStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[32m" // green for success
	case status >= 300 && status < 400:
		return "\033[34m" // blue for redirects
	case status >= 400 && status < 500:
		return "\033[33m" // yellow for client errors
	default:
		return "\033[31m" // red for server errors
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[34m" //blue
	case "POST":
		return "\033[32m" // green
	case "PUT":
		return "\033[33m" // yellow
	case "DELETE":
		return "\033[31m" // red
	default:
		return "\033[37m" // white
	}
}

package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

// TraceID returns a middleware that assigns a trace ID to every request, at the very front
// of the chain (input), so it's available to every downstream layer (logging, proxy, response).
// If the caller already sent X-Trace-Id (e.g. propagated from another upstream gateway), that
// value is kept — otherwise a new UUID is generated. The trace ID is echoed back to the client
// on X-Trace-Id so it can be correlated with logs/support tickets, and it's what
// proxy.injectCallerHeaders forwards upstream on the same header (output to service).
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request.Context(), helpers.KeyTraceID, traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Set("trace_id", traceID)
		c.Header("X-Trace-Id", traceID)

		c.Next()
	}
}

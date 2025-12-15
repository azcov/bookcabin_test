package httpz

import (
	"context"
	"net/http"
	"time"

	"github.com/azcov/bookcabin_test/pkg/consts"
	"github.com/azcov/bookcabin_test/pkg/errorz"
	"github.com/azcov/bookcabin_test/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID middleware adds or extracts an X-Request-ID header into the context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(consts.HeaderRequestID)
		if _, err := uuid.Parse(rid); err != nil {
			rid = uuid.New().String()
			c.Request.Header.Set(consts.HeaderRequestID, rid)
			c.Writer.Header().Set(consts.HeaderRequestID, rid)
		}
		c.Set(consts.HeaderRequestID, rid)

		// Update the request context so it is available via c.Request.Context()
		ctx := context.WithValue(c.Request.Context(), consts.HeaderRequestID, rid)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// Logger middleware logs the request path, method, status code and latency
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		rid := ""
		if v, ok := c.Get(consts.HeaderRequestID); ok {
			rid = v.(string)
		}
		logger.InfoContext(c.Request.Context(), "http_request", "method", c.Request.Method, "path", c.Request.URL.Path, "status", status, "latency_ms", latency.Milliseconds(), "rid", rid)
	}
}

// Recovery middleware catches panics and returns a standardized JSON error.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(c.Request.Context(), "panic_recovered", "err", r)
				we := &errorz.WrappedError{StatusCode: http.StatusInternalServerError, ErrCode: "internal_error", Msg: "internal server error"}
				c.Abort()
				JSONResponse(c, nil, we)
			}
		}()
		c.Next()
	}
}

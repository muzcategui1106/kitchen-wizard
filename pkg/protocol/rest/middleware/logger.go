package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ctxLogger struct{}

const (
	loggerKey = "logger"
)

// StructuredLogger logs request/response pair
func StructuredLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Do not log Kubernetes health check.
		// You can change this behavior as you wish.
		if ctx.Request.Header.Get("X-Liveness-Probe") == "Healthz" {
			ctx.Next()
			return
		}

		id := GetReqID(ctx)

		// Prepare fields to log
		var scheme string
		if ctx.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
		method := ctx.Request.Method
		remoteAddr := ctx.Request.RemoteAddr
		userAgent := ctx.Request.UserAgent()
		uri := strings.Join([]string{scheme, "://", ctx.Request.Host, ctx.Request.RequestURI}, "")

		logger = logger.WithOptions(zap.Fields(zap.String("request-id", id)))
		LoogerToContext(ctx, logger)

		// Log HTTP request
		logger.Debug("request started",
			zap.String("http-scheme", scheme),
			zap.String("http-method", method),
			zap.String("remote-addr", remoteAddr),
			zap.String("user-agent", userAgent),
			zap.String("uri", uri),
		)

		ctx.Next()

		// Log HTTP response
		logger.Debug("request completed",
			zap.String("http-scheme", scheme),
			zap.String("http-method", method),
			zap.String("remote-addr", remoteAddr),
			zap.String("user-agent", userAgent),
			zap.String("uri", uri),
			zap.Float64("latency", float64(time.Since(startTime).Nanoseconds())/1000000.0),
		)
	}
}

// LoogerToContext sets logger to context
func LoogerToContext(ctx *gin.Context, l *zap.Logger) {
	ctx.Set(loggerKey, l)
}

// LoggerFromContext returns logger from context
func LoggerFromContext(ctx *gin.Context) *zap.Logger {
	li, ok := ctx.Get(loggerKey)
	if !ok {
		return zap.L()
	}

	l, ok := li.(*zap.Logger)

	if ok {
		return l
	}

	return zap.L()
}

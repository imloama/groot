package groot

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

func TraceIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(TRACE_ID)
		if traceId == "" {
			traceId = xid.New().String()
		}
		ctx := context.WithValue(context.Background(), TRACE_ID, traceId)
		c.Set(TRACE_CTX, ctx)

		c.Next()
	}
}

func GetTraceId(c *gin.Context) (string, bool) {
	ctxval, exit := c.Get(TRACE_CTX)
	if !exit {
		return "", false
	}
	ctx := ctxval.(context.Context)
	traceId := ctx.Value(TRACE_ID)
	if traceId != nil {
		return traceId.(string), true
	}
	return "", false
}

func TDebug(c *gin.Context, msg string, fields ...zap.Field) {
	traceId, exit := GetTraceId(c)
	if exit {
		lg.Debug(msg, append(fields, zap.String(TRACE_ID, traceId))...)
		return
	}
	lg.Debug(msg, fields...)
}

func TError(c *gin.Context, msg string, fields ...zap.Field) {
	traceId, exit := GetTraceId(c)
	if exit {
		lg.Error(msg, append(fields, zap.String(TRACE_ID, traceId))...)
		return
	}
	lg.Error(msg, fields...)
}

func TInfo(c *gin.Context, msg string, fields ...zap.Field) {
	traceId, exit := GetTraceId(c)
	if exit {
		lg.Info(msg, append(fields, zap.String(TRACE_ID, traceId))...)
		return
	}
	lg.Info(msg, fields...)
}

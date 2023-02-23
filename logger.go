package groot

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxsize"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

var lg *zap.Logger

// InitLogger 初始化Logger
func InitLogger(cfg *LogConfig) (err error) {
	if cfg == nil || cfg.Filename == "" {
		cfg = &LogConfig{
			Level:      "debug",
			Filename:   "./logs/app.log",
			MaxSize:    100,
			MaxAge:     100,
			MaxBackups: 100,
		}
	} else {
		if cfg.Level == "" {
			cfg.Level = "debug"
		}
		if cfg.Filename == "" {
			cfg.Filename = "./logs/app.log"
		}
		if cfg.MaxSize <= 0 {
			cfg.MaxSize = 100
		}
		if cfg.MaxAge <= 0 {
			cfg.MaxAge = 100
		}
		if cfg.MaxBackups <= 0 {
			cfg.MaxBackups = 100
		}
	}
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writes := []zapcore.WriteSyncer{writeSyncer, zapcore.AddSync(os.Stdout)}
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic(err)
	}
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writes...), l)
	lg = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(lg)
	return
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func Debug(msg string, fields ...zap.Field) {
	Log().Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	lg.Error(msg, fields...)
}

func Log() *zap.Logger {
	return lg
}

// func (l *zap.Logger) GetCtx(ctx context.Context) *zap.Logger {
// 	log, ok := ctx.Value(l.opts.CtxKey).(*zap.Logger)
// 	if ok {
// 		return log
// 	}
// 	return l.Logger
// }

// func (l *zap.Logger) AddCtx(ctx context.Context, field ...zap.Field) (context.Context, *zap.Logger) {
// 	log := l.With(field...)
// 	ctx = context.WithValue(ctx, l.opts.CtxKey, log)
// 	return ctx, log
// }

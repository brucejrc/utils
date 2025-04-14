package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

type Logger interface {
	Debugw(msg string, kvs ...any)

	Infow(msg string, kvs ...any)

	Warnw(msg string, kvs ...any)

	Errorw(msg string, kvs ...any)

	Panicw(msg string, kvs ...any)

	Fatalw(msg string, kvs ...any)

	Sync()
}

var (
	mu  sync.Mutex
	std *zaplogger = New(NewOptions())
)

var _ Logger = (*zaplogger)(nil)

// zaplogger是Logger接口的具体实现，底层封装了zap.Logger
type zaplogger struct {
	z *zap.Logger
}

func Debugw(msg string, kvs ...any) {
	std.Debugw(msg, kvs...)
}

func (z *zaplogger) Debugw(msg string, kvs ...any) {
	z.z.Sugar().Debugw(msg, kvs...)
}

func Infow(msg string, kvs ...any) {
	std.Infow(msg, kvs...)
}

func (z *zaplogger) Infow(msg string, kvs ...any) {
	z.z.Sugar().Infow(msg, kvs...)
}

func Warnw(msg string, kvs ...any) {
	std.Warnw(msg, kvs...)
}

func (z *zaplogger) Warnw(msg string, kvs ...any) {
	z.z.Sugar().Warnw(msg, kvs...)
}

func Errorw(msg string, kvs ...any) {
	std.Errorw(msg, kvs...)
}

func (z *zaplogger) Errorw(msg string, kvs ...any) {
	z.z.Sugar().Errorw(msg, kvs...)
}

func Panicw(msg string, kvs ...any) {
	std.Panicw(msg, kvs...)
}

func (z *zaplogger) Panicw(msg string, kvs ...any) {
	z.z.Sugar().Panicw(msg, kvs...)
}

func Fatalw(msg string, kvs ...any) {
	std.Fatalw(msg, kvs...)
}

func (z *zaplogger) Fatalw(msg string, kvs ...any) {
	z.z.Sugar().Fatalw(msg, kvs...)
}

func Sync() {
	std.Sync()
}

func (z *zaplogger) Sync() {
	_ = z.z.Sync()
}

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

func New(opts *Options) *zaplogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.MessageKey = "message"

	encoderConfig.TimeKey = "timestamp"

	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-12 15:00:01.222"))
	}

	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	cfg := &zap.Config{
		DisableStacktrace: opts.DisableStacktrace,
		DisableCaller:     opts.DisableCaller,
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Encoding:          opts.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
	}

	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	zap.RedirectStdLog(z)

	return &zaplogger{z: z}
}

func W(ctx context.Context) Logger {
	return std.W(ctx)
}

func (z *zaplogger) W(ctx context.Context) Logger {
	lc := z.clone()
	// do append key and value from ctx to zaplogger
	// 定义一个映射，关联 context 提取函数和日志字段名。
	// contextExtractors := map[string]func(context.Context) string{
	//	known.XRequestID: contextx.RequestID, // 提取请求 ID
	//	known.XUserID:    contextx.UserID,    // 提取用户 ID
	// }

	// 遍历映射，从 context 中提取值并添加到日志中。
	// for fieldName, extractor := range contextExtractors {
	//	if val := extractor(ctx); val != "" {
	//		lc.z = lc.z.With(zap.String(fieldName, val))
	//	}
	//}
	return lc
}

func (z *zaplogger) clone() *zaplogger {
	newLogger := *z
	return &newLogger
}

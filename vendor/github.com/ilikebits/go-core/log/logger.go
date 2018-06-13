package log

import (
	"context"
	"runtime"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

func From(ctx context.Context) *Logger {
	reqLogger := zerolog.Ctx(ctx).With().Str("RequestID", middleware.GetReqID(ctx)).Logger()
	return &Logger{&reqLogger}
}

// TODO: alternatively, logging frame + 1 caller upwards? use pkg/file:line
func (l *Logger) Panic() *zerolog.Event {
	if dbg {
		return logWithCaller(l.Logger.Panic())
	}
	return l.Logger.Panic()
}

func (l *Logger) Fatal() *zerolog.Event {
	if dbg {
		return logWithCaller(l.Logger.Fatal())
	}
	return l.Logger.Fatal()
}

func (l *Logger) Error() *zerolog.Event {
	if dbg {
		return logWithCaller(l.Logger.Error())
	}
	return l.Logger.Error()
}

func (l *Logger) Warn() *zerolog.Event {
	if dbg {
		return logWithCaller(l.Logger.Warn())
	}
	return l.Logger.Warn()
}

func (l *Logger) Info() *zerolog.Event {
	if dbg {
		return logWithCaller(l.Logger.Info())
	}
	return l.Logger.Info()
}

func (l *Logger) Debug() *zerolog.Event {
	if dbg {
		return logWithCaller(l.Logger.Debug())
	}
	return l.Logger.Debug()
}

func logWithCaller(e *zerolog.Event) *zerolog.Event {
	_, file, line, ok := runtime.Caller(3) // Caller(3) logger.X -> log.X -> caller
	if !ok {
		panic("unable to get caller")
	}
	return e.Str("file", file).Int("line", line)
}

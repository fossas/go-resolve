// Package log implements application-level logging functions.
package log

import (
	"context"

	"github.com/apex/log"
)

const ContextLoggerKey = "go_resolve_logger"

func (l *Entry) Panic(msg string) {
}

func (l *Entry) Panicf(msg string, v ...interface{}) {

}

func FromContext(ctx context.Context) *log.Entry {
	return log.NewEntry(log.Log.(*log.Logger)).WithFields(log.Fields{
		"request-ID": ctx.Value("requestID"),
	})
}

func Panic(msg string) {

}

func Panicf(msg string, v ...interface{}) {

}

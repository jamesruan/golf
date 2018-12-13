package golf

import (
	"github.com/jamesruan/golf/event"
)

type Entry interface {
	Output(calldepth int, level event.Level, format string, args []interface{})
	OutputFatal(calldepth int, level event.Level, format string, args []interface{})
	OutputPanic(calldepth int, level event.Level, format string, args []interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(...event.Field) Entry
}

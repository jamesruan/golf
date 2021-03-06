package golf

import (
	"github.com/Workiva/go-datastructures/list"
	"github.com/jamesruan/golf/event"
	"os"
)

type SimpleEntry struct {
	fields  list.PersistentList
	handler event.Handler
}

func NewSimpleEntry(handler event.Handler) Entry {
	return &SimpleEntry{
		fields:  list.Empty,
		handler: handler,
	}
}

func (t *SimpleEntry) WithFields(fields ...event.Field) Entry {
	l := t.fields
	for _, field := range fields {
		l = l.Add(field)
	}
	return &SimpleEntry{
		fields:  l,
		handler: t.handler,
	}
}

func (t *SimpleEntry) Output(calldepth int, evel event.Level, format string, args []interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) OutputFatal(calldepth int, evel event.Level, format string, args []interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	os.Exit(-1)
}

func (t *SimpleEntry) OutputPanic(calldepth int, evel event.Level, format string, args []interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	panic("panic")
}

func (t *SimpleEntry) Debugf(format string, args ...interface{}) {
	t.Output(0, event.NOLEVEL, format, args)
}

func (t *SimpleEntry) Infof(format string, args ...interface{}) {
	t.Output(0, event.NOLEVEL, format, args)
}

func (t *SimpleEntry) Warnf(format string, args ...interface{}) {
	t.Output(0, event.NOLEVEL, format, args)
}

func (t *SimpleEntry) Errorf(format string, args ...interface{}) {
	t.Output(0, event.NOLEVEL, format, args)
}

func (t *SimpleEntry) Fatalf(format string, args ...interface{}) {
	t.OutputFatal(0, event.NOLEVEL, format, args)
}

func (t *SimpleEntry) Panicf(format string, args ...interface{}) {
	t.OutputPanic(0, event.NOLEVEL, format, args)
}

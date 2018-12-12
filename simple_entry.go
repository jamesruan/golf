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

func NewSimpleEntry(handler event.Handler) *SimpleEntry {
	return &SimpleEntry{
		fields:  list.Empty,
		handler: handler,
	}
}

func (t *SimpleEntry) WithFields(fields ...event.Field) *SimpleEntry {
	l := t.fields
	for _, field := range fields {
		l = l.Add(field)
	}
	return &SimpleEntry{
		fields:  l,
		handler: t.handler,
	}
}

func (t *SimpleEntry) Debugf(format string, args ...interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Infof(format string, args ...interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Warnf(format string, args ...interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Errorf(format string, args ...interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Fatalf(format string, args ...interface{}) {
	e := event.SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	os.Exit(-1)
}

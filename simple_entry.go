package golf

import (
	"github.com/Workiva/go-datastructures/list"
	"os"
)

type SimpleEntry struct {
	fields  list.PersistentList
	handler Handler
}

func NewSimpleEntry(handler Handler) *SimpleEntry {
	return &SimpleEntry{
		fields:  list.Empty,
		handler: handler,
	}
}

func (t *SimpleEntry) WithFields(fields ...EventField) *SimpleEntry {
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
	e := SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Infof(format string, args ...interface{}) {
	e := SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Warnf(format string, args ...interface{}) {
	e := SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Errorf(format string, args ...interface{}) {
	e := SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
}

func (t *SimpleEntry) Fatalf(format string, args ...interface{}) {
	e := SimpleEvent(format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	os.Exit(-1)
}

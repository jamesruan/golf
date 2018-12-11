package golf

import (
	"github.com/Workiva/go-datastructures/list"
	"os"
)

type Entry interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	WithFields(...EventField) Entry
}

type TopicEntry struct {
	topic   string
	fields  list.PersistentList
	handler Handler
}

func NewTopicEntry(topic string, handler Handler) *TopicEntry {
	return &TopicEntry{
		topic:   topic,
		fields:  list.Empty,
		handler: handler,
	}
}

func (t *TopicEntry) WithFields(fields ...EventField) *TopicEntry {
	l := t.fields
	for _, field := range fields {
		l = l.Add(field)
	}
	return &TopicEntry{
		topic:   t.topic,
		fields:  l,
		handler: t.handler,
	}
}

func (t *TopicEntry) Debugf(format string, args ...interface{}) {
	e := DefaultEvent(3, t.topic, DEBUG, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Infof(format string, args ...interface{}) {
	e := DefaultEvent(3, t.topic, INFO, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Warnf(format string, args ...interface{}) {
	e := DefaultEvent(3, t.topic, WARN, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Errorf(format string, args ...interface{}) {
	e := DefaultEvent(3, t.topic, ERROR, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Fatalf(format string, args ...interface{}) {
	e := DefaultEvent(3, t.topic, FATAL, format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	os.Exit(-1)
}

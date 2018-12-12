package golf

import (
	"github.com/Workiva/go-datastructures/list"
	"github.com/jamesruan/golf/event"
	"os"
)

type TopicEntry struct {
	topic   string
	fields  list.PersistentList
	handler event.Handler
}

func NewTopicEntry(topic string, handler event.Handler) *TopicEntry {
	return &TopicEntry{
		topic:   topic,
		fields:  list.Empty,
		handler: handler,
	}
}

func (t *TopicEntry) WithFields(fields ...event.Field) Entry {
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
	e := event.DefaultEvent(3, t.topic, event.DEBUG, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Infof(format string, args ...interface{}) {
	e := event.DefaultEvent(3, t.topic, event.INFO, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Warnf(format string, args ...interface{}) {
	e := event.DefaultEvent(3, t.topic, event.WARN, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Errorf(format string, args ...interface{}) {
	e := event.DefaultEvent(3, t.topic, event.ERROR, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) Fatalf(format string, args ...interface{}) {
	e := event.DefaultEvent(3, t.topic, event.FATAL, format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	os.Exit(-1)
}

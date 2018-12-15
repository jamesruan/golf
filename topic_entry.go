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

func NewTopicEntry(topic string, handler event.Handler) Entry {
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

func (t *TopicEntry) Output(calldepth int, level event.Level, format string, args []interface{}) {
	e := event.DefaultEvent(3+calldepth, t.topic, level, format, args, t.fields)
	t.handler.Handle(e)
}

func (t *TopicEntry) OutputFatal(calldepth int, level event.Level, format string, args []interface{}) {
	e := event.DefaultEvent(3+calldepth, t.topic, level, format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	os.Exit(-1)
}

func (t *TopicEntry) OutputPanic(calldepth int, level event.Level, format string, args []interface{}) {
	e := event.DefaultEvent(3+calldepth, t.topic, level, format, args, t.fields)
	t.handler.Handle(e)
	close(sinkCloseSignal)
	sinkWg.Wait()
	panic("panic")
}

func (t *TopicEntry) Debugf(format string, args ...interface{}) {
	t.Output(1, event.DEBUG, format, args)
}

func (t *TopicEntry) Infof(format string, args ...interface{}) {
	t.Output(1, event.INFO, format, args)
}

func (t *TopicEntry) Warnf(format string, args ...interface{}) {
	t.Output(1, event.WARN, format, args)
}

func (t *TopicEntry) Errorf(format string, args ...interface{}) {
	t.Output(1, event.ERROR, format, args)
}

func (t *TopicEntry) Fatalf(format string, args ...interface{}) {
	t.OutputFatal(1, event.FATAL, format, args)
}

func (t *TopicEntry) Panicf(format string, args ...interface{}) {
	t.OutputPanic(1, event.PANIC, format, args)
}

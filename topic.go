package golf

import (
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/processor"
)

func NewTopicLogHandler(topic string) LogHandler {
	return &topicLogHandler{
		topic: topic,
	}
}

type topicLogHandler struct {
	topic string
}

func (t topicLogHandler) Debugf(fmt string, args ...interface{}) {
	e := event.New(2, t.topic, event.DEBUG, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Infof(fmt string, args ...interface{}) {
	e := event.New(2, t.topic, event.INFO, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Logf(fmt string, args ...interface{}) {
	e := event.New(2, t.topic, event.LOG, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Warnf(fmt string, args ...interface{}) {
	e := event.New(2, t.topic, event.WARN, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Errorf(fmt string, args ...interface{}) {
	e := event.New(2, t.topic, event.ERROR, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Processor(next processor.P) processor.P {
	return processor.NewNamedP(t.topic, next)
}

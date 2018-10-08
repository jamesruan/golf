package golf

import (
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/processor"
)

// NewTopicLogHandler create a new logger that has embedded topic.
// It is by default connected to DiscardLoggerP
func NewTopicLogHandler(topic string) LogHandler {
	t := &topicLogHandler{
		topic: topic,
	}
	t.Processor(DiscardLoggerP)
	return t
}

type topicLogHandler struct {
	topic string
}

func (t topicLogHandler) Debugf(fmt string, args ...interface{}) {
	e := event.Default(3, t.topic, event.DEBUG, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Infof(fmt string, args ...interface{}) {
	e := event.Default(3, t.topic, event.INFO, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Logf(fmt string, args ...interface{}) {
	e := event.Default(3, t.topic, event.LOG, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Warnf(fmt string, args ...interface{}) {
	e := event.Default(3, t.topic, event.WARN, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Errorf(fmt string, args ...interface{}) {
	e := event.Default(3, t.topic, event.ERROR, fmt, args, nil)
	mainP.Process(e)
}

func (t topicLogHandler) Fatalf(fmt string, args ...interface{}) {
	e := event.Default(3, t.topic, event.ERROR, fmt, args, nil)
	mainP.Process(e)
	processor.Exit()
}

func (t topicLogHandler) Processor(next processor.P) processor.P {
	return processor.NewNamedP(t.topic, next)
}

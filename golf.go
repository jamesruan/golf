package golf

import (
	"context"
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/logger"
	"github.com/jamesruan/golf/processor"
)

var mainP = processor.NewTopicP("golf", context.Background())

var (
	DefaultLoggerP = processor.NewLoggerP("stderr", logger.DefaultStderrLogger)
	DiscardLoggerP = processor.NewLoggerP("discard", logger.DiscardLogger)
	DefaultP       = processor.NewLogLevelP(event.INFO).Either(DefaultLoggerP).Or(DiscardLoggerP)
)

var ()

func init() {
	defaultTopicP := NewTopicLogHandler("").Processor(DefaultP)
	RegisterTopicProcessor(defaultTopicP)
}

type LogHandler interface {
	Processor(processor.P) processor.P
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Logf(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
}

func Debugf(fmt string, args ...interface{}) {
	e := event.Default(3, "", event.DEBUG, fmt, args, nil)
	mainP.Process(e)
}

func Infof(fmt string, args ...interface{}) {
	e := event.Default(3, "", event.INFO, fmt, args, nil)
	mainP.Process(e)
}

func Logf(fmt string, args ...interface{}) {
	e := event.Default(3, "", event.LOG, fmt, args, nil)
	mainP.Process(e)
}

func Warnf(fmt string, args ...interface{}) {
	e := event.Default(3, "", event.WARN, fmt, args, nil)
	mainP.Process(e)
}

func Errorf(fmt string, args ...interface{}) {
	e := event.Default(3, "", event.ERROR, fmt, args, nil)
	mainP.Process(e)
}

func Printf(fmt string, args ...interface{}) {
	e := event.Simple(fmt, args, nil)
	mainP.Process(e)
}

// RegisterTopicProcessor register a processor with it's name the topic name
func RegisterTopicProcessor(p processor.P) {
	mainP.Set(p)
}

package golf

import (
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/logger"
	"github.com/jamesruan/golf/processor"
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	Debugf("%d", 0)
	Infof("%d", 1)
	Logf("%d", 2)
	Warnf("%d", 3)
	Errorf("%d", 4)

	frameLoggerP := processor.NewLoggerP("stderr", logger.NewConsoleLogger(os.Stderr, logger.LstdFlags|logger.Lframes))
	p := processor.NewLogLevelP(event.DEBUG).Either(frameLoggerP).Or(DiscardLoggerP)
	v := NewTopicLogHandler("mytopic")

	RegisterTopicProcessor(v.Processor(p))
	v.Debugf("%d", 0)
	v.Infof("%d", 1)
	v.Logf("%d", 2)
	v.Warnf("%d", 3)
	v.Errorf("%d", 4)
}

func BenchmarkDefaultStderr(b *testing.B) {
	RegisterTopicProcessor(NewTopicLogHandler("").Processor(DiscardLoggerP))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Infof("test")
		}
	})
}

func BenchmarkSimpleStderr(b *testing.B) {
	RegisterTopicProcessor(NewTopicLogHandler("").Processor(DiscardLoggerP))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Printf("test")
		}
	})
}

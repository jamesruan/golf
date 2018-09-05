package golf

import (
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/logger"
	"github.com/jamesruan/golf/processor"
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	Debugf("%d", 0)
	Infof("%d", 1)
	Logf("%d", 2)
	Warnf("%d", 3)
	Errorf("%d", 4)

	frameConsoleLogger := logger.NewConsoleLogger(os.Stderr, logger.LstdFlags|logger.Lframes)
	frameLoggerP := processor.NewLoggerP("stderr", frameConsoleLogger)
	p := processor.NewLogLevelP(event.DEBUG).Either(frameLoggerP).Or(DiscardLoggerP)

	v := NewTopicLogHandler("mytopic")

	RegisterTopicProcessor(v.Processor(p))
	v.Debugf("%d", 0)
	v.Infof("%d", 1)
	v.Logf("%d", 2)
	v.Warnf("%d", 3)
	v.Errorf("%d", 4)
	time.Sleep(1 * time.Second)
}

func BenchmarkDefault(b *testing.B) {
	RegisterTopicProcessor(NewTopicLogHandler("").Processor(DiscardLoggerP))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Infof("test")
		}
	})
}

func BenchmarkSimple(b *testing.B) {
	RegisterTopicProcessor(NewTopicLogHandler("").Processor(DiscardLoggerP))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Printf("test")
		}
	})
}

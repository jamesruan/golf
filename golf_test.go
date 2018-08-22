package golf

import (
	"testing"
)

func TestDefault(t *testing.T) {
	Debugf("%d", 0)
	Infof("%d", 1)
	Logf("%d", 2)
	Warnf("%d", 3)
	Errorf("%d", 4)

	v := NewTopicLogHandler("mytopic")
	p := v.Processor(DefaultP)
	RegisterTopicProcessor(p)
	v.Debugf("%d", 0)
	v.Infof("%d", 1)
	v.Logf("%d", 2)
	v.Warnf("%d", 3)
	v.Errorf("%d", 4)
}

func BenchmarkStderr(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Infof("test")
		}
	})
}

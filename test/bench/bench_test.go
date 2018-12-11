package bench

import (
	"github.com/jamesruan/golf"
	"github.com/jamesruan/golf/sinks/console"
	"testing"
)

func BenchmarkDefault(b *testing.B) {
	infohandler := golf.NewLevelHandler(golf.INFO, console.Discard, nil)
	broadcastor := golf.NewBroadcastHandler()
	broadcastor.AddHandler("1", infohandler)
	broadcastor.AddHandler("2", infohandler)
	broadcastor.AddHandler("3", infohandler)
	broadcastor.AddHandler("4", infohandler)
	test_entry := golf.NewTopicEntry("test", broadcastor)
	test_entry1 := test_entry.WithFields(golf.EventField{Name: "app", Value: "test"})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			test_entry1.Errorf("5")
		}
	})
}

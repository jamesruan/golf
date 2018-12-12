package bench

import (
	"github.com/jamesruan/golf"
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/formatter/text"
	"github.com/jamesruan/golf/handlers"
	"io/ioutil"
	"testing"
)

func BenchmarkDefault(b *testing.B) {
	discardTextHandler := golf.DefaultStreamSink(ioutil.Discard, text.Console)
	infohandler := handlers.NewLevel(event.INFO, discardTextHandler, nil)
	broadcastor := handlers.NewBroadcast()
	broadcastor.AddHandler("1", infohandler)
	broadcastor.AddHandler("2", infohandler)
	broadcastor.AddHandler("3", infohandler)
	broadcastor.AddHandler("4", infohandler)
	test_entry := golf.NewTopicEntry("test", broadcastor)
	test_entry1 := test_entry.WithFields(event.Field{Name: "app", Value: "test"})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			test_entry1.Errorf("5")
		}
	})
}

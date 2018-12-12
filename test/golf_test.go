package golf

import (
	"github.com/jamesruan/golf"
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/formatter/text"
	"github.com/jamesruan/golf/handlers"
	"os"
	"testing"
	"time"
)

func TestNewTopicEntry(t *testing.T) {

	broadcastor := handlers.NewBroadcast()
	callstackSink := golf.DefaultStreamSink(os.Stderr, text.New(text.Lframes))
	broadcastor.AddHandler("callstack", callstackSink)
	broadcastor.AddHandler("raw", golf.DefaultPlainSink)

	infohandler := handlers.NewLevel(event.INFO, golf.DefaultSink, nil)
	broadcastor.AddHandler("info", infohandler)

	test_entry := golf.NewTopicEntry("test", broadcastor)
	test_entry.Debugf("1")
	test_entry.Infof("2")
	test_entry.Warnf("3")
	test_entry1 := test_entry.WithFields(event.Field{Name: "app", Value: "test"})
	test_entry2 := test_entry1.WithFields(event.Field{Name: "module", Value: "testmodule"})
	test_entry.Errorf("4")
	time.Sleep(time.Second)
	test_entry2.Errorf("5")
	test_entry1.Fatalf("5")
}

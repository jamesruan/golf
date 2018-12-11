package golf

import (
	"github.com/jamesruan/golf"
	"github.com/jamesruan/golf/sinks/console"
	"testing"
	"time"
)

func TestNewTopicEntry(t *testing.T) {
	broadcastor := golf.NewBroadcastHandler()
	broadcastor.AddHandler("raw", console.Default)
	infohandler := golf.NewLevelHandler(golf.INFO, console.Default, nil)
	broadcastor.AddHandler("info", infohandler)
	test_entry := golf.NewTopicEntry("test", broadcastor)
	test_entry.Debugf("1")
	test_entry.Infof("2")
	test_entry.Warnf("3")
	test_entry1 := test_entry.WithFields(golf.EventField{Name: "app", Value: "test"})
	test_entry2 := test_entry1.WithFields(golf.EventField{Name: "module", Value: "testmodule"})
	test_entry.Errorf("4")
	time.Sleep(time.Second)
	test_entry2.Errorf("5")
	test_entry1.Fatalf("5")
}
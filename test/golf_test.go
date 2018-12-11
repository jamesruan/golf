package golf

import (
	"github.com/jamesruan/golf"
	"github.com/jamesruan/golf/sinks/console"
	"testing"
)

func TestNewTopicEntry(t *testing.T) {
	handler := golf.NewSinkHandler(console.Default)
	entry := golf.NewTopicEntry("test", handler)
	entry.Infof("test")
	entry1 := entry.WithFields(golf.EventField{Name: "app", Value: "test"})
	entry.Errorf("test fatal")
	entry1.Fatalf("test fatal")
}

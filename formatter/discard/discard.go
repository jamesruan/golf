package discard

import (
	"github.com/jamesruan/golf/event"
)

var (
	Default = &DiscardFormatter{}
)

type DiscardFormatter struct{}

func (f *DiscardFormatter) Format(e *event.Event) []byte {
	return nil
}

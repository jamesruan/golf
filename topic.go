package golf

import (
	"github.com/jamesruan/golf/processor"
)

func NewTopicProcessor(name string, next processor.P) processor.P {
	return processor.NewNamedP(name, next)
}

package golf

import (
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/logger"
	"io"
)

// DirectSink sinks event with Formatter and sink the bytes direct into the Logger
type DirectSink struct {
	sinkBase
	output    io.Writer
	formatter event.Formatter
}

func NewDirectSink(output io.Writer, formatter event.Formatter) *DirectSink {
	ds := &DirectSink{
		output:    output,
		formatter: formatter,
	}
	ds.register()

	go func() {
		defer ds.deregister()
		for {
			select {
			case <-SinkRotateSignal:
				if r, ok := output.(logger.RotateLogger); ok {
					err := r.Rotate()
					if err != nil {
						panic(err)
					}
				}
			case <-ds.closing():
				return
			}
		}
	}()
	return ds
}

// Handle blocks until the event can be written
func (l DirectSink) Handle(e *event.Event) {
	l.output.Write(l.formatter.Format(e))
}

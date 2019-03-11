package golf

import (
	"bufio"
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/logger"
	"io"
)

// StreamSink sinks event with Formatter and sink the bytes into bufferred stream
type StreamSink struct {
	sinkBase
	bufout    *bufio.Writer
	formatter event.Formatter
}

func DefaultStreamSink(output io.Writer, formatter event.Formatter) *StreamSink {
	return NewStreamSink(output, formatter)
}

func NewStreamSink(output io.Writer, formatter event.Formatter) *StreamSink {
	bufout := bufio.NewWriter(output)
	ss := &StreamSink{
		bufout:    bufout,
		formatter: formatter,
	}
	ss.register()
	go func() {
		defer ss.deregister()
		for {
			select {
			case <-SinkRotateSignal:
				if r, ok := output.(logger.RotateLogger); ok {
					bufout.Flush()
					err := r.Rotate()
					if err != nil {
						panic(err)
					}
				}
			case <-ss.closing():
				bufout.Flush()
				return
			}
		}
	}()
	return ss
}

func (l *StreamSink) Handle(e *event.Event) {
	l.bufout.Write(l.formatter.Format(e))
}

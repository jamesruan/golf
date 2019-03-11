package golf

import (
	"bufio"
	"github.com/jamesruan/golf/event"
	"github.com/jamesruan/golf/logger"
	"time"
)

// StreamSink sinks event with Formatter and sink the bytes into bufferred stream
type StreamSink struct {
	sinkBase
	bufout    *bufio.Writer
	formatter event.Formatter
	queue     chan []byte
}

func DefaultStreamSink(output logger.StreamLogger, formatter event.Formatter) *StreamSink {
	return NewStreamSink(output, formatter, 16, 100*time.Millisecond)
}

func NewStreamSink(output logger.StreamLogger, formatter event.Formatter, bufferSize uint, maxDelay time.Duration) *StreamSink {
	queue := make(chan []byte, bufferSize)
	bufout := bufio.NewWriter(output)
	ss := &StreamSink{
		bufout:    bufout,
		formatter: formatter,
		queue:     queue,
	}
	ss.register()
	go func() {
		defer ss.deregister()
		ticker := time.NewTicker(maxDelay)
		ch := ticker.C
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
				for {
					select {
					case b := <-queue:
						bufout.Write(b)
					default:
						bufout.Flush()
						ticker.Stop()
						return
					}
				}
			case b := <-queue:
				bufout.Write(b)
				ch = ticker.C
			case <-ch:
				// flush every 100 microsecond if not flushed before
				bufout.Flush()
				ch = nil
			}
		}
	}()
	return ss
}

func (l StreamSink) Handle(e *event.Event) {
	b := l.formatter.Format(e)
	l.queue <- b
}

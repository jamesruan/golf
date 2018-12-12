package golf

import (
	"bufio"
	"github.com/jamesruan/golf/event"
	"io"
	"time"
)

// StreamSink sinks event with Formatter and sink the bytes into bufferred stream
type StreamSink struct {
	bufout    *bufio.Writer
	formatter event.Formatter
	queue     chan []byte
}

func DefaultStreamSink(output io.Writer, formatter event.Formatter) *StreamSink {
	return NewStreamSink(output, formatter, 16, 100*time.Millisecond)
}

func NewStreamSink(output io.Writer, formatter event.Formatter, bufferSize uint, maxDelay time.Duration) *StreamSink {
	queue := make(chan []byte, bufferSize)
	bufout := bufio.NewWriter(output)
	sinkWg.Add(1)
	go func() {
		defer sinkWg.Done()
		ticker := time.NewTicker(maxDelay)
		ch := ticker.C
		for {
			select {
			case <-sinkCloseSignal:
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
	return &StreamSink{
		bufout:    bufout,
		formatter: formatter,
		queue:     queue,
	}
}

func (l StreamSink) Handle(e *event.Event) {
	b := l.formatter.Format(e)
	l.queue <- b
}

package golf

import (
	"io"
	"github.com/jamesruan/cqueue"
)

func NewConsoleLogger(out io.Writer) *ConsoleLogger {
	return &ConsoleLogger{
		out: out,
		q: cqueue.New(),
	}
}

type ConsoleLogger struct {
	out io.Writer
	q cqueue.Queue
}

func (l *ConsoleLogger) Log(e *Event) {
}

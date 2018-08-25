package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jamesruan/golf/event"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync"
)

var (
	DefaultStderrLogger = NewConsoleLogger(os.Stderr, LstdFlags)
	DiscardLogger       = NewConsoleLogger(ioutil.Discard, LstdFlags)
)

func NewConsoleLogger(out io.Writer, flags ConsoleLoggerFlags) *ConsoleLogger {
	return &ConsoleLogger{
		out:    out,
		bufout: bufio.NewWriter(out),
		flags:  flags,
		queue:  make(chan *event.Event, 32),
	}
}

type ConsoleLogger struct {
	out    io.Writer
	bufout *bufio.Writer
	flags  ConsoleLoggerFlags
	queue  chan *event.Event
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 128)
		return bytes.NewBuffer(buf)
	},
}

func (l ConsoleLogger) Queue() <-chan *event.Event {
	return l.queue
}

func (l ConsoleLogger) Enqueue(e *event.Event) {
	l.queue <- e
}

func (l *ConsoleLogger) fmt(b *bytes.Buffer, e *event.Event) {
	simple := e.Level == event.NOLEVEL
	if !simple {
		if l.flags&CLcolor != 0 {
			fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m ", l.levelColor(e.Level), e.Level)
		} else {
			fmt.Fprintf(b, "[%s] ", e.Level)
		}
	}

	if l.flags&Ldatetime != 0 {
		var s string
		if l.flags&Lmicroseconds != 0 {
			s = e.Time.Format("2006/01/02 15:04:05.000 ")
		} else {
			s = e.Time.Format("2006/01/02 15:04:05 ")
		}
		b.WriteString(s)
	}

	if len(e.Topic) > 0 {
		fmt.Fprintf(b, "\x1b[92m%s\x1b[0m ", e.Topic)
	}

	if !simple && l.flags&Lframes == 0 {
		if len(e.Pc) > 0 {
			pc := e.Pc[0]
			f := runtime.FuncForPC(pc)
			file, line := f.FileLine(pc)
			if l.flags&Llongfile != 0 {
				fmt.Fprintf(b, "%s:%d: ", file, line)
			} else {
				fmt.Fprintf(b, "%s:%d: ", path.Base(file), line)
			}
		} else {
			fmt.Fprint(b, "???:0: ")
		}
	}

	if len(e.Args) > 0 {
		fmt.Fprintf(b, e.Fmt, e.Args...)
	} else {
		b.WriteString(e.Fmt)
	}

	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			fmt.Fprintf(b, "%s=%v ", k, v)
		}
	}

	if !simple && l.flags&Lframes != 0 && len(e.Pc) > 0 {
		frames := runtime.CallersFrames(e.Pc)
		for {
			f, more := frames.Next()
			fmt.Fprintf(b, "\n\t%s (in %s:%d)", f.Function, path.Base(f.File), f.Line)
			if !more {
				break
			}
		}
	}

	b.WriteByte(byte('\n'))
}

func (l *ConsoleLogger) Log(e *event.Event) {
	b := bufPool.Get().(*bytes.Buffer)
	l.fmt(b, e)
	l.bufout.Write(b.Bytes())
	b.Reset()
	bufPool.Put(b)
	if len(l.queue) == 0 {
		l.bufout.Flush()
	}
}

func (ConsoleLogger) levelColor(l event.Level) int {
	switch l {
	case event.DEBUG:
		return 37 //white
	case event.INFO:
		return 34 //blue
	case event.WARN:
		return 33 //yellow
	case event.LOG:
		return 32 //green
	case event.ERROR:
		return 31 //red
	default:
		return 0 // nocolor
	}
}

type ConsoleLoggerFlags uint32

const (
	Llongfile     = 1 << iota           // full file name and line number: /a/b/c/d.go:23
	Ldatetime                           // the date in the local time zone: 2009/01/23
	Lmicroseconds                       // microsecond resolution: 01:23:23.123123.  assumes Ldatatime.
	Lframes                             // display calling stack frames
	CLcolor                             // colorize
	LstdFlags     = CLcolor | Ldatetime // initial values for the standard logger
)

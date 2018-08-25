package logger

import (
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
		out:   out,
		flags: flags,
	}
}

type ConsoleLogger struct {
	out   io.Writer
	flags ConsoleLoggerFlags
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 128)
		return bytes.NewBuffer(buf)
	},
}

func (l *ConsoleLogger) fmt(b *bytes.Buffer, e *event.Event) {
	if l.flags&CLcolor != 0 {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m ", l.levelColor(e.Level), e.Level)
	} else {
		fmt.Fprintf(b, "[%s] ", e.Level)
	}

	if l.flags&Ldatetime != 0 {
		Y, M, D := e.Time.Date()
		h, m, s := e.Time.Clock()
		fmt.Fprintf(b, "%04d/%02d/%02d %02d:%02d:%02d ", Y, M, D, h, m, s)
		if l.flags&Lmicroseconds != 0 {
			n := e.Time.Nanosecond()
			fmt.Fprintf(b, ".%03d", n)
		}
	}

	if len(e.Topic) > 0 {
		fmt.Fprintf(b, "\x1b[97m%s\x1b[0m ", e.Topic)
	}

	if l.flags&Lframes == 0 {
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

	fmt.Fprintf(b, e.Fmt, e.Args...)

	if l.flags&Lframes != 0 && len(e.Pc) > 0 {
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
	if l.out == ioutil.Discard {
		return
	}
	b := bufPool.Get().(*bytes.Buffer)
	l.fmt(b, e)
	l.out.Write(b.Bytes())
	b.Reset()
	bufPool.Put(b)
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

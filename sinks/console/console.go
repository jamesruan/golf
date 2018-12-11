package console

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jamesruan/golf"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

var (
	DefaultMaxFlushDelay = 50 * time.Millisecond
	DefaultHiRes         = golf.NewSinkHandler(New(os.Stderr, Ldatetime|Lmicroseconds, 16, DefaultMaxFlushDelay))
	DefaultPlainText     = golf.NewSinkHandler(New(os.Stderr, Ldatetime, 16, DefaultMaxFlushDelay))
	Default              = golf.NewSinkHandler(New(os.Stderr, LstdFlags, 16, DefaultMaxFlushDelay))
	Discard              = golf.NewSinkHandler(New(ioutil.Discard, LstdFlags, 16, DefaultMaxFlushDelay))
)

type ConsoleSink struct {
	out    io.Writer
	bufout *bufio.Writer
	flags  ConsoleSinkFlags
	queue  chan *golf.Event
	done   chan struct{}
}

func New(out io.Writer, flags ConsoleSinkFlags, bufferSize int, maxDelay time.Duration) *ConsoleSink {
	queue := make(chan *golf.Event, bufferSize)
	done := make(chan struct{})
	sink := &ConsoleSink{
		out:    out,
		bufout: bufio.NewWriter(out),
		flags:  flags,
		queue:  queue,
		done:   done,
	}
	go func() {
		ticker := time.NewTicker(maxDelay)
		ch := ticker.C
		for {
			select {
			case e, ok := <-queue:
				if !ok {
					sink.bufout.Flush()
					ticker.Stop()
					close(done)
					return
				}
				sink.log(e)
				ch = ticker.C
			case <-ch:
				// flush every 100 microsecond if not flushed before
				sink.bufout.Flush()
				ch = nil
			}
		}
	}()
	return sink
}

func (l ConsoleSink) Name() string {
	return "console"
}

func (l ConsoleSink) Close() {
	close(l.queue)
	<-l.done
}

func (l ConsoleSink) Handle(e *golf.Event) {
	l.queue <- e
}

func (l *ConsoleSink) fmt(b *bytes.Buffer, e *golf.Event) {
	if l.flags&Ldatetime != 0 {
		var s string
		if l.flags&Lmicroseconds != 0 {
			s = e.Time.Format("2006/01/02 15:04:05.000 ")
		} else {
			s = e.Time.Format("2006/01/02 15:04:05 ")
		}
		b.WriteString(s)
	}

	simple := e.Level == golf.NOLEVEL
	if !simple {
		if l.flags&CLcolor != 0 {
			fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m ", l.levelColor(e.Level), e.Level)
		} else {
			fmt.Fprintf(b, "[%s] ", e.Level)
		}
	}

	if len(e.Topic) > 0 {
		fmt.Fprintf(b, "\x1b[92m%s\x1b[0m ", e.Topic)
	}

	if !simple && l.flags&Lframes == 0 {
		if len(e.Pc) > 0 {
			frames := runtime.CallersFrames(e.Pc)
			f, more := frames.Next()
			for f.File == "<autogenerated>" && more {
				f, more = frames.Next()
			}
			if l.flags&Llongfile != 0 {
				fmt.Fprintf(b, "%s:%d: ", f.File, f.Line)
			} else {
				fmt.Fprintf(b, "%s:%d: ", path.Base(f.File), f.Line)
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

	if l := e.Fields.Length(); l > 0 {
		var i uint
		for i = 0; i < l; i++ {
			field, _ := e.Fields.Get(i)
			f := field.(golf.EventField)
			fmt.Fprintf(b, " %s=%v", f.Name, f.Value)
		}
	}

	if !simple && l.flags&Lframes != 0 && len(e.Pc) > 0 {
		frames := runtime.CallersFrames(e.Pc)
		for {
			f, more := frames.Next()
			for f.File == "<autogenerated>" && more {
				f, more = frames.Next()
			}
			fmt.Fprintf(b, "\n\t%s (in %s:%d)", f.Function, path.Base(f.File), f.Line)
			if !more {
				break
			}
		}
	}

	b.WriteByte(byte('\n'))
}

func (l *ConsoleSink) log(e *golf.Event) {
	b := bufPool.Get().(*bytes.Buffer)
	l.fmt(b, e)
	l.bufout.Write(b.Bytes())
	b.Reset()
	bufPool.Put(b)
}

func (ConsoleSink) levelColor(l golf.Level) int {
	switch l {
	case golf.DEBUG:
		return 37 //white
	case golf.INFO:
		return 34 //blue
	case golf.WARN:
		return 33 //yellow
	case golf.ERROR:
		fallthrough
	case golf.FATAL:
		return 31 //red
	default:
		return 0 // nocolor
	}
}

type ConsoleSinkFlags uint32

const (
	Llongfile     = 1 << iota           // full file name and line number: /a/b/c/d.go:23
	Ldatetime                           // the date in the local time zone: 2009/01/23
	Lmicroseconds                       // microsecond resolution: 01:23:23.123123.  assumes Ldatatime.
	Lframes                             // display calling stack frames
	CLcolor                             // colorize
	LstdFlags     = CLcolor | Ldatetime // initial values for the standard logger
)

var bufPool = &sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 256)
		return bytes.NewBuffer(buf)
	},
}

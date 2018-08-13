package golf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
)

var stderrLogger = NewConsoleLogger(os.Stderr, LstdFlags)

func init() {
	stderrLogger.SetFilter(WithLevel(LOG))
}

func NewConsoleLogger(out io.Writer, flags ConsoleLoggerFlags) *ConsoleLogger {

	return &ConsoleLogger{
		out:   out,
		flags: flags,
	}
}

type ConsoleLogger struct {
	out    io.Writer
	flags  ConsoleLoggerFlags
	filter EventFilter
}

func (l *ConsoleLogger) SetFilter(f EventFilter) {
	l.filter = f
}

func (l *ConsoleLogger) Log(e *Event) {
	if !l.filter.FilterEvent(e) {
		return
	}

	b := new(bytes.Buffer)

	if l.flags&CLcolor != 0 {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m ", l.levelColor(e.Level), e.Level)
	} else {
		fmt.Fprintf(b, "[%s] ", e.Level)
	}

	if l.flags&Ldatetime != 0 {
		Y, M, D := e.Time.Date()
		h, m, s := e.Time.Clock()
		fmt.Fprintf(b, "%04d/%02d/%02d %02d:%02d:%02d", Y, M, D, h, m, s)
		if l.flags&Lmicroseconds != 0 {
			n := e.Time.Nanosecond()
			fmt.Fprintf(b, ".%03d", n)
		}
	}

	fmt.Fprintf(b, "\x1b[97m%s\x1b[0m ", e.Topic)

	if l.flags&Llongfile != 0 {
		fmt.Fprintf(b, "%s:%d: ", e.File, e.Line)
	} else {
		fmt.Fprintf(b, "%s:%d: ", path.Base(e.File), e.Line)
	}

	fmt.Fprintf(b, e.Fmt+"\n", e.Args...)

	l.out.Write(b.Bytes())
}

func (ConsoleLogger) levelColor(l Level) int {
	switch l {
	case DEBUG:
		return 37 //white
	case INFO:
		return 34 //blue
	case WARN:
		return 33 //yellow
	case LOG:
		return 32 //green
	case ERROR:
		return 31 //red
	default:
		return 0 // nocolor
	}
}

type ConsoleLoggerFlags uint32

const (
	Lshortfile    = 0                                // final file name element and line number: d.go:23
	Llongfile     = 1 << iota                        // full file name and line number: /a/b/c/d.go:23
	Ldatetime                                        // the date in the local time zone: 2009/01/23
	Lmicroseconds                                    // microsecond resolution: 01:23:23.123123.  assumes Ldatatime.
	CLcolor                                          // colorize
	LstdFlags     = CLcolor | Ldatetime | Lshortfile // initial values for the standard logger
)

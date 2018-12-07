package processor

import (
	"github.com/jamesruan/golf/event"
	"strings"
	"sync"
)

type logLevelP struct {
	level event.Level
	state string

	ch_either   chan P
	ch_or       chan P
	ch_flushing chan struct{}
	ch_event    chan *event.Event
	ch_stopped  chan struct{}
	either      P
	or          P
}

func (l logLevelP) Name() string {
	var sb strings.Builder
	sb.WriteString("LogLevel")
	sb.WriteString(l.level.String())
	return sb.String()
}

func (l logLevelP) Process(e *event.Event) {
	l.ch_event <- e
}

func (l *logLevelP) Either(p P) EitherP {
	l.ch_either <- p
	return l
}

func (l *logLevelP) Or(p P) EitherP {
	l.ch_or <- p
	return l
}

func (l logLevelP) Judge(e *event.Event) bool {
	return e.Level >= l.level
}

func (l logLevelP) process(e *event.Event) {
	if l.Judge(e) {
		if l.either != nil {
			l.either.Process(e)
		}
	} else {
		if l.or != nil {
			l.or.Process(e)
		}
	}
}

func (l logLevelP) Stopped() <-chan struct{} {
	return l.ch_stopped
}

func (l *logLevelP) loop(stop <-chan struct{}) {
	defer close(l.ch_stopped)

	for {
		switch l.state {
		case "INIT":
			if l.either != nil || l.or != nil {
				l.state = "ACTIVE"
			}
			select {
			case p := <-l.ch_either:
				l.either = p
			case p := <-l.ch_or:
				l.or = p
			case <-stop:
				l.state = "STOPPED"
			}
		case "ACTIVE":
			select {
			case p := <-l.ch_either:
				l.either = p
			case p := <-l.ch_or:
				l.or = p
			case e := <-l.ch_event:
				l.process(e)
			case <-l.ch_flushing:
				l.flush()
			case <-stop:
				l.state = "STOPPING"
			}
		case "STOPPING":
			select {
			case e := <-l.ch_event:
				l.process(e)
			default:
				l.state = "STOPPED"
			}
		case "STOPPED":
			l.flush()
			return
		}
	}
}
func (l *logLevelP) flush() {
	print("FLUSH LOGLEVEL\n")
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		if l.either != nil {
			l.either.Flush()
		}
		wg.Done()
	}()
	go func() {
		if l.or != nil {
			l.or.Flush()
		}
		wg.Done()
	}()
	wg.Wait()
}

func (l *logLevelP) Start(stop <-chan struct{}) P {
	go l.loop(stop)
	return l
}

func (l logLevelP) Flush() {
	l.ch_flushing <- struct{}{}
}

// NewLogLevelP returns a processor that process event.Level >= lvl in Either branch.
func NewLogLevelP(lvl event.Level) EitherP {
	return &logLevelP{
		level:       lvl,
		state:       "INIT",
		ch_event:    make(chan *event.Event),
		ch_either:   make(chan P, 1),
		ch_or:       make(chan P, 1),
		ch_flushing: make(chan struct{}),
		ch_stopped:  make(chan struct{}),
	}
}

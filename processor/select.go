package processor

import (
	"fmt"
	"github.com/jamesruan/golf/event"
)

type selectFunc = func(map[string]P, *event.Event) ([]P, bool)

type selectP struct {
	name        string
	processors  map[string]P
	f           selectFunc
	ch_started  chan struct{}
	ch_event    chan *event.Event
	ch_command  chan selectPCmd
	ch_flushing chan struct{}
	ch_stopping chan struct{}
	ch_stopped  chan struct{}
}

type selectPCmd struct {
	p P      // for add
	n string // for delete
}

func (t selectP) Name() string {
	return t.name
}

func (t selectP) Process(e *event.Event) {
	<-t.ch_started
	select {
	case <-t.ch_stopping:
		return
	default:
		t.ch_event <- e
	}
}

func (t selectP) Set(p P) {
	<-t.ch_started
	t.ch_command <- selectPCmd{p: p}
}

func (t selectP) Unset(name string) {
	<-t.ch_started
	t.ch_command <- selectPCmd{n: name}
}

func (t selectP) Select(e *event.Event) ([]P, bool) {
	return t.f(t.processors, e)
}

func (t selectP) process(e *event.Event, flushing bool) {
	if ps, ok := t.Select(e); ok {
		for i, _ := range ps {
			ps[i].Process(e)
			if flushing {
				ps[i].Flush()
			}
		}
	}
}

func (t selectP) Stopped() <-chan struct{} {
	return t.ch_stopped
}

func (t selectP) Stop() {
	select {
	case <-t.ch_stopped:
		panic(fmt.Sprintf("golf: Stop an already Stopped P: %s", t.Name()))
	default:
	}
	close(t.ch_stopping)
}

func (t *selectP) Start(stop <-chan struct{}) P {
	select {
	case <-t.ch_started:
		panic(fmt.Sprintf("golf: Start an already started P: %s", t.Name()))
	default:
	}
	go func() {
		close(t.ch_started)
		flushing := false
		parent_stop := stop
		for {
			select {
			case cmd := <-t.ch_command:
				if cmd.p != nil {
					t.processors[cmd.p.Name()] = cmd.p
				} else {
					delete(t.processors, cmd.n)
				}
			case <-t.ch_stopping:
				// flush pending event
				for {
					select {
					case e := <-t.ch_event:
						t.process(e, true)
					default:
						close(t.ch_stopped)
						return
					}
				}
			default:
			}
			select {
			case e := <-t.ch_event:
				t.process(e, flushing)
				if len(t.ch_event) == 0 {
					flushing = false
				}
			case <-parent_stop:
				t.Stop()
				parent_stop = nil
			case <-t.ch_flushing:
				flushing = true
			}
		}
	}()

	return t
}

func (t selectP) Flush() {
}

func makeSelectP(name string, f selectFunc) selectP {
	t := selectP{
		name:        name,
		processors:  make(map[string]P),
		f:           f,
		ch_started:  make(chan struct{}),
		ch_event:    make(chan *event.Event, 16),
		ch_command:  make(chan selectPCmd, 1),
		ch_flushing: make(chan struct{}),
		ch_stopping: make(chan struct{}),
		ch_stopped:  make(chan struct{}),
	}

	return t
}

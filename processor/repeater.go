package processor

import (
	"github.com/jamesruan/golf/event"
)

// EventSelectFunc select where the event should be sent to from targets
type EventSelectFunc func(processors map[string]P, e *event.Event)

// Repeater process the message by its internal selector and send the event to selected targets
type Repeater struct {
	name     string
	selector *eventSelector
	done     chan struct{}
}

type eventSelector struct {
	Processors  map[string]P
	fun         EventSelectFunc
	ch_fun      chan EventSelectFunc
	ch_event    chan *event.Event
	ch_selected chan []string
	ch_addP     chan P
	ch_delP     chan string
}

func newEventSelector(done chan struct{}) *eventSelector {
	s := &eventSelector{
		Processors: make(map[string]P),
		fun:        defaultEventSelectFunc,
		ch_fun:     make(chan EventSelectFunc),
		ch_event:   make(chan *event.Event, 1),
		ch_addP:    make(chan P),
		ch_delP:    make(chan string),
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case f := <-s.ch_fun:
				if f != nil {
					s.fun = f
				} else {
					s.fun = defaultEventSelectFunc
				}
			case p := <-s.ch_addP:
				s.Processors[p.Name()] = p
			case n := <-s.ch_delP:
				delete(s.Processors, n)
			case e := <-s.ch_event:
				s.fun(s.Processors, e)
			}
		}
	}()

	return s
}

func defaultEventSelectFunc(processors map[string]P, e *event.Event) {
	for _, p := range processors {
		p.Process(e)
	}
}

func (r Repeater) Process(e *event.Event) {
	r.selector.ch_event <- e
}

func (r Repeater) Name() string {
	return r.name
}

func (r Repeater) Set(p P) {
	r.selector.ch_addP <- p
}

func (r Repeater) Unset(name string) {
	r.selector.ch_delP <- name
}

func (r Repeater) Selector(fun EventSelectFunc) {
	r.selector.ch_fun <- fun
}

func (r Repeater) Close() {
	close(r.done)
}

func NewRepeater(name string) *Repeater {
	done := make(chan struct{})
	return &Repeater{
		name:     name,
		selector: newEventSelector(done),
		done:     done,
	}
}

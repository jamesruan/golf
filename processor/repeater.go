package processor

import (
	"github.com/jamesruan/golf/event"
	"sync"
)

// EventSelectFunc select where the event should be sent to from targets
type EventSelectFunc func(targets []string, e *event.Event) (selected []string)

// Repeater process the message by its internal selector and send the event to selected targets
type Repeater struct {
	name       string
	processors map[string]P
	selector   *eventSelector
	lock       *sync.RWMutex // make update processors and notifying selector atomic
	done       chan struct{}
}

type eventSelector struct {
	fun         EventSelectFunc
	targets     []string
	ch_fun      chan EventSelectFunc
	ch_targets  chan []string
	ch_event    chan *event.Event
	ch_selected chan []string
}

func newEventSelector(done chan struct{}) *eventSelector {
	s := &eventSelector{
		fun:         defaultEventSelectFunc,
		targets:     []string{},
		ch_fun:      make(chan EventSelectFunc),
		ch_targets:  make(chan []string),
		ch_event:    make(chan *event.Event),
		ch_selected: make(chan []string),
	}

	go func() {
		var selected []string
		var event chan *event.Event = s.ch_event
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
			case t := <-s.ch_targets:
				s.targets = t
			case e := <-event:
				event = nil
				selected = s.fun(s.targets, e)
			case s.ch_selected <- selected:
				event = s.ch_event
			}
		}
	}()

	return s
}

func defaultEventSelectFunc(targets []string, e *event.Event) []string {
	return targets
}

func (r Repeater) Process(e *event.Event) {
	r.selector.ch_event <- e
	selected := <-r.selector.ch_selected

	for _, s := range selected {
		r.lock.RLock()
		p, ok := r.processors[s]
		r.lock.RUnlock()
		if ok {
			p.Process(e)
		}
	}
}

func (r Repeater) Name() string {
	return r.name
}

func (r Repeater) Set(p P) {
	name := p.Name()
	targets := make([]string, 0, len(r.processors)+1)
	r.lock.Lock()
	r.processors[name] = p
	for k, _ := range r.processors {
		targets = append(targets, k)
	}
	r.selector.ch_targets <- targets
	r.lock.Unlock()
}

func (r Repeater) Unset(name string) {
	r.lock.Lock()
	targets := make([]string, 0, len(r.processors)-1)
	delete(r.processors, name)
	for k, _ := range r.processors {
		targets = append(targets, k)
	}
	r.selector.ch_targets <- targets
	r.lock.Unlock()
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
		name:       name,
		processors: make(map[string]P),
		selector:   newEventSelector(done),
		done:       done,
		lock:       new(sync.RWMutex),
	}
}
package processor

import (
	"github.com/jamesruan/golf/event"
	"sync"
)

type selectFunc = func(map[string]P, *event.Event) ([]P, bool)

type selectP struct {
	name        string
	processors  map[string]P
	f           selectFunc
	state       string
	ch_flushing chan struct{}
	ch_event    chan *event.Event
	ch_add      chan P
	ch_delete   chan string
	ch_stopped  chan struct{}
}

func (t selectP) Name() string {
	return t.name
}

func (t selectP) Process(e *event.Event) {
	t.ch_event <- e
}

func (t selectP) Set(p P) {
	t.ch_add <- p
}

func (t selectP) Unset(name string) {
	t.ch_delete <- name
}

func (t selectP) Select(e *event.Event) ([]P, bool) {
	return t.f(t.processors, e)
}

func (t selectP) process(e *event.Event) {
	if ps, ok := t.Select(e); ok {
		for i, _ := range ps {
			ps[i].Process(e)
		}
	}
}

func (t selectP) Stopped() <-chan struct{} {
	return t.ch_stopped
}

func (t *selectP) loop(stop <-chan struct{}) {
	defer close(t.ch_stopped)

	for {
		switch t.state {
		case "INIT":
			select {
			case p := <-t.ch_add:
				t.processors[p.Name()] = p
				t.state = "ACTIVE"
			case <-stop:
				t.state = "STOPPED"
			}
		case "ACTIVE":
			select {
			case p := <-t.ch_add:
				t.processors[p.Name()] = p
				t.state = "ACTIVE"
			case name := <-t.ch_delete:
				delete(t.processors, name)
			case e := <-t.ch_event:
				t.process(e)
			case <-t.ch_flushing:
				t.flush()
			case <-stop:
				t.state = "STOPPING"
			}
		case "STOPPING":
			select {
			case e := <-t.ch_event:
				t.process(e)
			default:
				t.state = "STOPPED"
			}
		case "STOPPED":
			t.flush()
			return
		}
	}
}

func (t *selectP) flush() {
	print("FLUSH SELECT\n")
	wg := new(sync.WaitGroup)
	for n, _ := range t.processors {
		wg.Add(1)
		go func(i string) {
			t.processors[i].Flush()
			wg.Done()
		}(n)
	}
	wg.Wait()
}

func (t *selectP) Start(stop <-chan struct{}) P {
	go t.loop(stop)

	return t
}

func (t selectP) Flush() {
	t.ch_flushing <- struct{}{}
}

func makeSelectP(name string, f selectFunc) selectP {
	t := selectP{
		name:        name,
		processors:  make(map[string]P),
		f:           f,
		state:       "INIT",
		ch_event:    make(chan *event.Event, 16),
		ch_flushing: make(chan struct{}),
		ch_add:      make(chan P),
		ch_delete:   make(chan string),
		ch_stopped:  make(chan struct{}),
	}

	return t
}

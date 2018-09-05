package processor

import (
	"context"
	"github.com/jamesruan/golf/event"
)

type selectFunc = func(map[string]P, *event.Event) ([]P, bool)

type selectP struct {
	name       string
	processors map[string]P
	f          selectFunc
	ch_event   chan *event.Event
	ch_addP    chan P
	ch_delP    chan string
	ctx        context.Context
	cancel     context.CancelFunc
}

func (t selectP) Name() string {
	return t.name
}

func (t selectP) Process(e *event.Event) {
	t.ch_event <- e
}

func (t selectP) Set(p P) {
	t.ch_addP <- p
}

func (t selectP) Unset(name string) {
	t.ch_delP <- name
}

func (t selectP) Select(e *event.Event) ([]P, bool) {
	return t.f(t.processors, e)
}

func (t selectP) Close() {
	t.cancel()
}

func makeSelectP(name string, ctx context.Context, f selectFunc) selectP {
	nctx, cancel := context.WithCancel(ctx)
	t := selectP{
		name:       name,
		processors: make(map[string]P),
		f:          f,
		ch_event:   make(chan *event.Event),
		ch_addP:    make(chan P),
		ch_delP:    make(chan string),
		ctx:        nctx,
		cancel:     cancel,
	}
	go func() {
		for {
			select {
			case <-nctx.Done():
				return
			case p := <-t.ch_addP:
				t.processors[p.Name()] = p
			case n := <-t.ch_delP:
				delete(t.processors, n)
			case e := <-t.ch_event:
				ps, ok := t.Select(e)
				if ok {
					for _, p := range ps {
						p.Process(e)
					}
				}
			}
		}
	}()

	return t
}

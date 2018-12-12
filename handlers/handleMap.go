package handlers

import (
	"github.com/jamesruan/golf/event"
	"sync"
)

type handlerMap sync.Map

func (m *handlerMap) Range(f func(key string, value event.Handler) bool) {
	(*sync.Map)(m).Range(func(k interface{}, v interface{}) bool {
		return f(k.(string), v.(event.Handler))
	})
}

func (m *handlerMap) Store(n string, h event.Handler) {
	(*sync.Map)(m).Store(n, h)
}

func (m *handlerMap) Load(n string) event.Handler {
	h, ok := (*sync.Map)(m).Load(n)
	if ok {
		return h.(event.Handler)
	} else {
		return nil
	}
}

func (m *handlerMap) Delete(n string) {
	(*sync.Map)(m).Delete(n)
}

func (m *handlerMap) Get(name string) event.Handler {
	h, ok := (*sync.Map)(m).Load(name)
	if ok && h != nil {
		return h.(event.Handler)
	}
	return nil
}

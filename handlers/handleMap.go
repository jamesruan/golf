package handlers

import (
	"github.com/jamesruan/golf"
	"sync"
)

type handlerMap sync.Map

func (m *handlerMap) Range(f func(key string, value golf.Handler) bool) {
	(*sync.Map)(m).Range(func(k interface{}, v interface{}) bool {
		return f(k.(string), v.(golf.Handler))
	})
}

func (m *handlerMap) Store(n string, h golf.Handler) {
	(*sync.Map)(m).Store(n, h)
}

func (m *handlerMap) Load(n string) golf.Handler {
	h, ok := (*sync.Map)(m).Load(n)
	if ok {
		return h.(golf.Handler)
	} else {
		return nil
	}
}

func (m *handlerMap) Delete(n string) {
	(*sync.Map)(m).Delete(n)
}

func (m *handlerMap) Get(name string) golf.Handler {
	h, ok := (*sync.Map)(m).Load(name)
	if ok && h != nil {
		return h.(golf.Handler)
	}
	return nil
}

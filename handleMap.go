package golf

import (
	"sync"
)

type handlerMap sync.Map

func (m *handlerMap) Range(f func(key string, value Handler) bool) {
	(*sync.Map)(m).Range(func(k interface{}, v interface{}) bool {
		return f(k.(string), v.(Handler))
	})
}

func (m *handlerMap) Store(n string, h Handler) {
	(*sync.Map)(m).Store(n, h)
}

func (m *handlerMap) Load(n string) Handler {
	h, ok := (*sync.Map)(m).Load(n)
	if ok {
		return h.(Handler)
	} else {
		return nil
	}
}

func (m *handlerMap) Delete(n string) {
	(*sync.Map)(m).Delete(n)
}

func (m *handlerMap) Get(name string) Handler {
	h, ok := (*sync.Map)(m).Load(name)
	if ok && h != nil {
		return h.(Handler)
	}
	return nil
}
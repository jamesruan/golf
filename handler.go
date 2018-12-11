package golf

type mapHandler struct {
	name string
	smap *handlerMap
}

func (h mapHandler) Name() string {
	return h.name
}

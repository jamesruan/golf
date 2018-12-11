package golf

// Handler handles Event.
type Handler interface {
	Handle(*Event)
}

type mapHandler struct {
	smap *handlerMap
}

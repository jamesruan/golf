package golf

type BroadcastHandler mapHandler

// NewBroadcastHandler creates a BroadcastHandler.
//
// BroadcastHandler handles event to all handlers registered.
// The order of calling registered handlers is random.
func NewBroadcastHandler() *BroadcastHandler {
	return &BroadcastHandler{
		name: "_broadcast",
		smap: new(handlerMap),
	}
}

func (h BroadcastHandler) Name() string {
	return h.name
}

func (h *BroadcastHandler) Handle(e *Event) {
	(*mapHandler)(h).smap.Range(func(name string, handler Handler) bool {
		handler.Handle(e)
		return true
	})
}

//AddHandler adds handler 'v' identified as 'id'
func (h *BroadcastHandler) AddHandler(id string, v Handler) {
	(*mapHandler)(h).smap.Store(id, v)
}

//DeleteHandler deletes handler identified as 'id'
func (h *BroadcastHandler) DeleteHandler(id string) {
	(*mapHandler)(h).smap.Delete(id)
}

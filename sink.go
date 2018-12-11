package golf

// SinkHandler handles event with sink and close it when necessary.
type SinkHandler struct {
	sink Sink
}

func (h *SinkHandler) Handle(e *Event) {
	h.sink.Handle(e)
}

func NewSinkHandler(sink Sink) *SinkHandler {
	sinkWg.Add(1)
	go func() {
		<-sinkCloseSignal
		sink.Close()
		sinkWg.Done()
	}()
	return &SinkHandler{sink}
}

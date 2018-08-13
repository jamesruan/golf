package golf

type levelFilter struct {
	level Level
}

func (l levelFilter) FilterEvent(e *Event) bool {
	return e.Level >= l.level
}

type EventFilter interface {
	FilterEvent(*Event) bool
}

func WithLevel(l Level) EventFilter {
	return levelFilter{l}
}

func SetFilter(topic string, f EventFilter) {
	h := getTopicHandler(topic)
	h.setFilter(f)
}

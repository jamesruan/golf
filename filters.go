package golf

type levelFilter struct {
	level Level
}

func (l levelFilter) FilterEvent(e *Event) bool {
	return e.Level >= l.level
}

func WithLevel(level Level) Filter {
	return levelFilter{level}
}

type topicFilter struct {
	topic string
}

func (l topicFilter) FilterEvent(e *Event) bool {
	return e.Topic == l.topic
}

func WithTopic(topic string) Filter {
	return topicFilter{topic}
}

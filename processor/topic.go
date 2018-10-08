package processor

import (
	"github.com/jamesruan/golf/event"
)

type topicP struct {
	selectP
}

// NewTopicP returns a processor that handle the event with a processor whose name is identical to event's topic.
func NewTopicP(name string) SelectP {
	t := makeSelectP(name, func(ps map[string]P, e *event.Event) ([]P, bool) {
		p, ok := ps[e.Topic]
		return []P{p}, ok
	})

	return &topicP{
		selectP: t,
	}
}
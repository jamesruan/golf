package processor

import (
	"context"
	"github.com/jamesruan/golf/event"
)

type repeaterP struct {
	selectP
}

// NewRepeaterP returns a processor that repeats event to each of its processors in random order.
func NewRepeaterP(name string, ctx context.Context) SelectP {
	t := makeSelectP(name, ctx, func(ps map[string]P, e *event.Event) ([]P, bool) {
		rps := make([]P, 0, len(ps))
		for _, p := range ps {
			rps = append(rps, p)
		}
		return rps, true
	})

	return &repeaterP{
		selectP: t,
	}
}

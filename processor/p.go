package processor

import (
	"context"
	"github.com/jamesruan/golf/event"
)

// P is a basic combinator for event processing
type P interface {
	Name() string
	Process(*event.Event)
}

// ResidentP is a P that running in a separate go routing thus needs to be started.
// It ends when End() is called or when ctx cancelled.
type ResidentP interface {
	P
	Context() context.Context
	Start(ctx context.Context) P
	End()
}

// EitherP is a P that choose down stream P by Judge()
type EitherP interface {
	ResidentP
	Judge(e *event.Event) bool
	Either(P) EitherP // Set down stream processors for Judge() returning true
	Or(P) EitherP     // Set down stream processors for Judge() returning false
}

// Select is a P that choose down stream P by Select()
type SelectP interface {
	ResidentP
	Select(*event.Event) ([]P, bool)
	Unset(name string)
	Set(P)
}

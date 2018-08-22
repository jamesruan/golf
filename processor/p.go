package processor

import (
	"github.com/jamesruan/golf/event"
)

type P interface {
	Name() string
	Process(*event.Event)
}

type EitherP interface {
	P
	Judge(e *event.Event) bool
	Either(P) EitherP // Set down stream processors for Judge() returning true
	Or(P) EitherP     // Set down stream processors for Judge() returning false
}

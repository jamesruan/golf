package processor

type namedP struct {
	P
	name string
}

func (p namedP) Name() string {
	return p.name
}

// NewNamedP warp a processor with new name.
func NewNamedP(name string, next P) P {
	return &namedP{
		name: name,
		P:    next,
	}
}

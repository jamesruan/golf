package processor

type NamedP struct {
	P
	name string
}

func (p NamedP) Name() string {
	return p.name
}

func NewNamedP(name string, next P) *NamedP {
	return &NamedP{
		name: name,
		P:    next,
	}
}

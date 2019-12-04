package agent

type Signal struct {
	ID uint64
}

func NewSignal() *Signal {
	a := &Signal{}
	return a
}

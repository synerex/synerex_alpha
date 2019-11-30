package agent

type Car struct {
	ID uint64
}

func NewCar() *Car {
	c := &Car{}
	return c
}

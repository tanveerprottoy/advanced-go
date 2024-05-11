package structs

const (
	defaultClass        = "sedan"
	defaultIsMod        = false
	defaultSuperCharger = ""
)

type Car struct {
	class        string
	mod          bool
	superCharger string
}

type Option func(*Car)

func NewCar(opts ...Option) *Car {
	h := &Car{
		class:        defaultClass,
		mod:          defaultIsMod,
		superCharger: defaultSuperCharger,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func WithClass(s string) Option {
	return func(c *Car) {
		c.class = s
	}
}

func WithMod(b bool) Option {
	return func(c *Car) {
		c.mod = b
	}
}

func WithSuperCharger(s string) Option {
	return func(c *Car) {
		c.superCharger = s
	}
}

func (c Car) Class() string {
	return c.class
}

func (c Car) Mod() bool {
	return c.mod
}

func (c Car) SuperCharger() string {
	return c.superCharger
}

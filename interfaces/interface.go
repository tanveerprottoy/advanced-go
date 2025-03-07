package interfaces

import "log"

type Doer interface {
	Do()
}

type Implementer struct {
	Err error
}

func (i *Implementer) Error() string {
	return "an error occurred on implementer"
}

func (i Implementer) Do() {
	log.Println("Do() called")
}

var _ Doer = Implementer{} // Verify that Implementer implements Doer.

var _ error = (*Implementer)(nil) // Verify that *Implementer implements error.

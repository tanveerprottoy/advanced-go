package validator

// Validator is an interface that defines the Validate method.
type Validator interface {
	Validate() error
}

// example struct
type User struct {
	Name  string `validate:"required",min=1,max=50`
	Email string `validate:"required",email`
	Age   int    `validate:"required",min=18,max=100`
}

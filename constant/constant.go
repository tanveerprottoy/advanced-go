package constant

/*
Go’s iota identifier is used in const declarations to simplify definitions of
incrementing numbers. Because it can be used in expressions, it provides a
generality beyond that of simple enumerations.

The value of iota is reset to 0 whenever the reserved word const appears
in the source (i.e. each const block) and incremented by one after each
ConstSpec e.g. each Line. This can be combined with the constant shorthand
(leaving out everything after the constant name) to very concisely
define related constants.
*/

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

type Weekday int

const (
	Saturday Weekday = iota
	Sunday
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
)

type Color int

const (
	Red Color = iota
	Green
	Blue
)

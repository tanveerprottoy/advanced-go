package constraint

import "log"

// From Go 1.25 CORE TYPE will be only applicable to type
// parameters (generics)

// For instance, interface{ ~[]int } has a core type ([]int), but
// the Constraint interface below does not have a core type. To
// make things more complicated, when it comes to channel
// operations and certain built-in calls (append, copy) the below
// definition of core types is too restrictive. The actual rules
// have adjustments that allow for differing channel directions
// and type sets containing both []byte and string types.

// There are various problems with this approach:

// Because the definition of core type must lead to sound type
// rules for different language features, it is overly restrictive
// for specific operations. For instance, the Go 1.24 rules for
// slice expressions do rely on core types, and as a consequence
// slicing an operand of type S constrained by Constraint is not
// permitted, even though it could be valid.

// When trying to understand a specific language feature, one may
// have to learn the intricacies of core types even when
// considering non-generic code. Again, for slice expressions,
// the language spec talks about the core type of the sliced
// operand, rather than just stating that the operand must be
// an array, slice, or string. The latter is more direct, simpler,
// and clearer, and doesn’t require knowing another concept that
// may be irrelevant in the concrete case.
// Because the notion of core types exists, the rules for index
// expressions, and len and cap (and others), which all eschew
// core types, appear as exceptions in the language rather than
// the norm. In turn, core types cause proposals such as issue
// #48522 which would permit a selector x.f to access a field f
// shared by all elements of x’s type set, to appear to add more
// exceptions to the language. Without core types, that feature
// becomes a natural and useful consequence of the ordinary rules
// for non-generic field access.

// Constraint is a generic constraint interface.
type Constraint interface {
	// the ~ operator, when used in type constraints within
	// generics, signifies an "approximation" of a type. It's
	// a way to express that you're not just looking for a
	// specific type, but also any type whose underlying type
	// is that specific type.
	~[]byte | ~string

	Hash() uint64
}

// the rules for index expressions state that (among other things)
// for an operand a of type parameter type P:
// The index expression a[x] must be valid for values of all types
// in P’s type set. The element types of all types in P’s type set
// must be identical. (In this context, the element type of a
// string type is byte.) These rules make it possible to index the
// generic variable s below
func at[T Constraint](s T, i int) byte {
	return s[i]
}

type MyString string

func (MyString) Hash() uint64 {
	panic("unimplemented")
}

func Executer() {
	log.Println(at(MyString("Hello"), 0))
}

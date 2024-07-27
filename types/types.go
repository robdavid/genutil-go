package types

import "golang.org/x/exp/constraints"

// Scalar numeric types with ordering
type Ranged interface {
	constraints.Float | constraints.Integer
}

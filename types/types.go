package types

import "golang.org/x/exp/constraints"

// Scalar numeric type constraint
type Real interface {
	constraints.Float | constraints.Integer
}

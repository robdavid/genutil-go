// Package yamlv3 provides specialized wrapper types for opt.Opt[T] which can be
// used with gopkg.in/yaml.v3 for marshaling and unmarshaling.
package yamlv3

import (
	"github.com/robdavid/genutil-go/opt"
	"gopkg.in/yaml.v3"
)

// Opt is an alias for [opt.Opt][T].
type Opt[T any] = opt.Opt[T]

// MutOpt is an alias for [opt.MutOpt][T].
type MutOpt[T any] = opt.MutOpt[T]

// OptVal is an alias for [opt.Val][T].
type OptVal[T any] = opt.Val[T]

// OptRef is an alias for [opt.Ref][T].
type OptRef[T any] = opt.Ref[T]

// Val wraps a standard optional value (like opt.Val[T]). This wrapper implements
// the full Opt[T] interface and provides specialized unmarshaling logic tailored
// for YAML v3 structs, allowing seamless use in struct tags.
type Val[T any] struct {
	OptVal[T]
}

// Ref wraps a reference optional value (like opt.Ref[T]). This wrapper also implements
// the full Opt[T] interface and provides specialized unmarshaling logic tailored
// for YAML v3 structs.
type Ref[T any] struct {
	OptRef[T]
}

// Value creates a Val[T] instance from an already existing value of type T.
func Value[T any](value T) Val[T] {
	return Val[T]{opt.Value(value)}
}

// Empty creates an empty (nil/zero) Val[T].
func Empty[T any]() Val[T] {
	return Val[T]{opt.Empty[T]()}
}

// Reference creates a Ref[T] instance from an existing pointer of type *T.
func Reference[T any](reference *T) Ref[T] {
	return Ref[T]{opt.Reference(reference)}
}

// EmptyRef creates an empty (nil/zero) Ref[T].
func EmptyRef[T any]() Ref[T] {
	return Ref[T]{opt.EmptyRef[T]()}
}

// UnmarshalYAML implements yaml.Unmarshaler for Val[T]. It decodes a YAML node
// into the underlying type T and wraps it in a new Value wrapper instance.
func (v *Val[T]) UnmarshalYAML(node *yaml.Node) error {
	var value T
	if err := node.Decode(&value); err != nil {
		return err
	}
	*v = Value(value)
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for Ref[T]. It decodes a YAML node
// into the underlying type T and wraps it in a new Reference wrapper instance.
func (v *Ref[T]) UnmarshalYAML(node *yaml.Node) error {
	var value T
	if err := node.Decode(&value); err != nil {
		return err
	}
	*v = Reference(&value)
	return nil
}

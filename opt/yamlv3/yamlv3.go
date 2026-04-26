package yamlv3

import (
	"github.com/robdavid/genutil-go/opt"
	"gopkg.in/yaml.v3"
)

type Option[T any] = opt.Option[T]

type Val[T any] struct {
	opt.Val[T]
}

type Ref[T any] struct {
	opt.Ref[T]
}

func Value[T any](value T) Val[T] {
	return Val[T]{opt.Value(value)}
}

func Empty[T any]() Val[T] {
	return Val[T]{opt.Empty[T]()}
}

func Reference[T any](reference *T) Ref[T] {
	return Ref[T]{opt.Reference(reference)}
}

func EmptyRef[T any]() Ref[T] {
	return Ref[T]{opt.EmptyRef[T]()}
}

func Equal[T comparable](o1 Option[T], o2 Option[T]) bool {
	return opt.Equal(o1, o2)
}

func (v *Val[T]) UnmarshalYAML(node *yaml.Node) error {
	var value T
	if err := node.Decode(&value); err != nil {
		return err
	}
	*v = Value(value)
	return nil
}

func (v *Ref[T]) UnmarshalYAML(node *yaml.Node) error {
	var value T
	if err := node.Decode(&value); err != nil {
		return err
	}
	*v = Reference(&value)
	return nil
}

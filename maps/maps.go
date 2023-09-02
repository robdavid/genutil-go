package maps

import (
	"errors"
	"fmt"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/option"
	"github.com/robdavid/genutil-go/slices"
	"github.com/robdavid/genutil-go/tuple"
)

// ErrPathConflict is an error constant that indicates that a nested map key path is being
// treated as both a value and a map.
var ErrPathConflict = errors.New("conflict between object and non-object types")

// ErrKeyError is an error constant that indicates that a key was not found.
var ErrKeyError = errors.New("key not found in map")

// PathConflict is an error type that indicates that a nested map key path is being
// treated as both a value and a map. It wraps ErrPathConflict and so
//
//	errors.Is(NewPathConflict(path), maps.ErrPathConflict)
//
// will return true.
type PathConflict[K comparable] []K

func (pc PathConflict[K]) Error() string {
	return fmt.Sprintf("%s at key path %v", ErrPathConflict, []K(pc))
}

func (pc PathConflict[K]) Unwrap() error {
	return ErrPathConflict
}

func NewPathConflict[K comparable](path []K) PathConflict[K] {
	return PathConflict[K](path)
}

// PathNotFound is an error type that indicates that a nested map does not
// contain the key path specified. It wraps ErrKeyError and so
//
//	errors.Is(NewPathNotFound(path), maps.ErrKeyError)
//
// will return true.
type PathNotFound[K comparable] []K

func (pnf PathNotFound[K]) Error() string {
	return fmt.Sprintf("%s at key path %v", ErrKeyError, []K(pnf))
}

func (pnf PathNotFound[K]) Unwrap() error {
	return ErrKeyError
}

func NewPathNotFound[K comparable](path []K) PathNotFound[K] {
	return PathNotFound[K](path)
}

// PutPath puts a value in a (possibly) nested map of maps. It mutates a map
// instance that maps from any comparable key to any value. The type of thew value
// can therefore encompass other map types. The path argument is a list of keys,
// each one representing a key at consecutive levels in the map. Intermediate maps
// are created as necessary. It is an error to attempt to replace an existing map
// with another value, or to replace an existing non-map value with a map.
//
//		m := make(map[string]any)
//		PutPath([]string{"a", "b"}, 123, m)
//		// m contains map[string]any{ "a": map[string]any{"b": 123} }
//	    err := PutPath([]string{"a"}, 456, m) // err != nil
func PutPath[K comparable, T any](path []K, value T, top map[K]T) error {
	m := top
	for i, s := range path {
		if i == len(path)-1 {
			// Reached final item in path
			if n, ok := m[s]; ok {
				if _, ok := any(n).(map[K]T); ok {
					return NewPathConflict(path)
				}
			}
			m[s] = value
		} else {
			if n, ok := m[s]; ok {
				if nm, okm := any(n).(map[K]T); okm {
					m = nm
				} else {
					return NewPathConflict(path[:i+1])
				}
			} else {
				n := any(make(map[K]T))
				m[s] = n.(T)
				m = n.(map[K]T)
			}
		}
	}
	return nil
}

func DeletePath[K comparable, T any](path []K, top map[K]T) (result T, ok bool, err error) {
	m := top
	result = any(top).(T)
	parents := make([]map[K]T, 0, len(path))
	parents = append(parents, top)
	for i, s := range path {
		if i == len(path)-1 {
			if result, ok = m[s]; ok {
				delete(m, s)
				for j := len(parents) - 1; j > 0; j-- {
					if len(parents[j]) == 0 {
						delete(parents[j-1], path[j-1])
					}
				}
			}
			break
		} else {
			var v T
			if v, ok = m[s]; ok {
				if m, ok = any(v).(map[K]T); ok {
					parents = append(parents, m)
				}
			}
			if !ok {
				break
			}
		}
	}
	return
}

func GetPath[K comparable, T any](path []K, top map[K]T) (result T, err error) {
	m := top
	result = any(top).(T)
	for i, s := range path {
		var ok bool
		if i == len(path)-1 {
			result, ok = m[s]
		} else {
			m, ok = any(m[s]).(map[K]T)
		}
		if !ok {
			err = NewPathNotFound(path[:i])
			return
		}
	}
	return
}

// Returns the keys of a map as a slice. The order of
// the keys is undefined.
func Keys[K comparable, T any](m map[K]T) []K {
	result := make([]K, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i++
	}
	return result
}

// Returns the values of a map as a slice. The order of the values
// is undefined.
func Values[K comparable, T any](m map[K]T) []T {
	result := make([]T, len(m))
	i := 0
	for _, v := range m {
		result[i] = v
		i++
	}
	return result
}

// Returns keys and values as a slice of 2-tuples. The order of the
// items is undefined
func Items[K comparable, T any](m map[K]T) []tuple.Tuple2[K, T] {
	result := make([]tuple.Tuple2[K, T], len(m))
	i := 0
	for k, v := range m {
		result[i] = tuple.Of2(k, v)
		i++
	}
	return result
}

// Returns keys and values as a slice of 2-tuples, sorted in key order
func SortedItems[K slices.Sortable, T any](m map[K]T) []tuple.Tuple2[K, T] {
	result := Items(m)
	slices.SortUsing(result, func(i1, i2 tuple.Tuple2[K, T]) bool { return i1.First < i2.First })
	return result
}

// Returns the keys of a map as a slice. The keys are sorted in their
// natural order, as defined by the < operator.
func SortedKeys[K slices.Sortable, T any](m map[K]T) []K {
	keys := Keys(m)
	slices.Sort(keys)
	return keys
}

// Returns the values of a map as a slice, sorted in the order
// of the associated key.
func SortedValuesByKey[K slices.Sortable, T any](m map[K]T) []T {
	return slices.Map(SortedKeys(m), AsFunc(m))
}

// Generate a function equivalent to a map, mapping keys to values.
func AsFunc[K comparable, T any](m map[K]T) func(K) T {
	return func(k K) (v T) {
		return m[k]
	}
}

// Finds a key value pair in a map which satisfies the predicate p which can match
// against both key and value. The return value is an option that either contains a
// 2-tuple of a matching key and value, or is empty. If there are multiple matching
// key value pairs, then which of those are returned is indeterminate.
func FindUsing[K comparable, T any](m map[K]T, p func(K, T) bool) option.Option[tuple.Tuple2[K, T]] {
	for k, v := range m {
		if p(k, v) {
			return option.Value(tuple.Of2(k, v))
		}
	}
	return option.Empty[tuple.Tuple2[K, T]]()
}

// Finds a key value pair in a map which satisfies the predicate p which can match
// against both key and value. The key and value are passed to the predicate function
// by reference. The return value is an option that either contains a 2-tuple of
// references to a matching key and value, or is empty. If there are multiple matching
// key value pairs, then which of those are returned is indeterminate.
func FindUsingRef[K comparable, T any](m map[K]T, p func(*K, *T) bool) option.Option[tuple.Tuple2[*K, *T]] {
	for k, v := range m {
		if p(&k, &v) {
			return option.Value(tuple.Of2(&k, &v))
		}
	}
	return option.Empty[tuple.Tuple2[*K, *T]]()
}

// Returns an iterator over the keys of a map.
func IterKeys[K comparable, T any](m map[K]T) iterator.Iterator[K] {
	return iterator.Generate(func(c iterator.Consumer[K]) {
		for k := range m {
			c.Yield(k)
		}
	})
}

// Returns an iterator over the values of a map.
func IterValues[K comparable, T any](m map[K]T) iterator.Iterator[T] {
	return iterator.Generate(func(c iterator.Consumer[T]) {
		for _, v := range m {
			c.Yield(v)
		}
	})
}

// Returns an iterator over the keys and values of a map, returning each pair
// as 2-tuple
func IterItems[K comparable, T any](m map[K]T) iterator.Iterator[tuple.Tuple2[K, T]] {
	return iterator.Generate(func(c iterator.Consumer[tuple.Tuple2[K, T]]) {
		for k, v := range m {
			c.Yield(tuple.Of2(k, v))
		}
	})
}

package maps

import (
	"errors"
	"fmt"

	"github.com/robdavid/genutil-go/slices"
)

var ErrPathConflict = errors.New("conflict between object and non-object types")
var ErrKeyError = errors.New("key not found in map")

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

func PutPath[K comparable, T any](path []K, value T, top map[K]T) error {
	m := top
	for i, s := range path {
		if i == len(path)-1 {
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

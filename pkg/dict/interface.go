// Dictionary implementation by Calder Lund

package dict

import (
	"tuple"
)


type DictInterface interface {
	/* Return the number of elements in dict. */
	Length() int
	/* Return the next key in dict. */
	Iterate() <-chan interface{}

	/* Remove key from the dict. */
	Remove(interface{}) error
	/* Returns the value at given key. */
	Get(interface{}) (interface{}, error)
	/* Sets the value at given key to given value. */
	Set(interface{}, interface{}) error
	/* Update the dict, adding elements from the other dict. Old values are replaced with new. */
	Combine(DictInterface) error
	/* Pop and return an arbitrary key from the dict. */
	PopKey() (interface{}, error)
	/* Pop and return an arbitrary value from the dict. */
	PopValue() (interface{}, error)
	/* Pop and return an arbitrary item from the dict. */
	Pop() (interface{}, interface{}, error)
	/* Clear all elements from the dict. */
	Clear() error

	/* Test for membership in the dict. */
	Contains(interface{}) (bool, error)
	/* Return true if the dict has all elements in common with the other dict. */
	Equals(DictInterface) (bool, error)

	/* Returns a tuple of keys */
	Keys() (tuple.TupleInterface, error)
	/* Returns a tuple of values. */
	Values() (tuple.TupleInterface, error)
	/* Returns a tuple of key/value pairs. */
	Items() (tuple.TupleInterface, error)

	/* Creates a copy of the current DictInterface */
	Copy() (DictInterface, error)

	/* Returns a string representation of the dict. */
	String() string

	/* Initializes the dict. */
	Init()
}

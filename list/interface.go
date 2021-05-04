package list

import "github.com/dynago/dg/internal/iterable"

// ListInterface is the interface which defines whether a struct is a list or not.
type ListInterface interface {
	/* Return the number of elements in list. */
	Length() int
	/* Return the next value in list. */
	Iterate() <-chan interface{}

	/* Test for membership in the list. */
	Contains(interface{}) (bool, error)
	/* Return true if the list has all elements in common with the other list. */
	Equals(ListInterface) (bool, error)

	/* Return concatenation of two lists together. */
	Concatenate(ListInterface) (ListInterface, error)
	/* Return list repeated n times. */
	Multiply(int) (ListInterface, error)
	/* Return reversed list. */
	Reverse() (ListInterface, error)

	/* Returns the value at index. */
	Get(int) (interface{}, error)
	/* Returns the a list of values given range. */
	Range(int, int) (ListInterface, error)
	/* Return first index of value. Returns -1 if not found. */
	Index(interface{}) (int, error)
	/* Return count of value. */
	Count(interface{}) (int, error)

	/* Inserts the value at index. */
	Insert(int, interface{}) error
	/* Sets values in given range to the values in the iterable. */
	Set(int, int, iterable.Iterable) error
	/* Remove first occurrence of the element from the list. */
	Remove(interface{}) error
	/* Removes range from the list. Takes two parameters: start (int), end (int: optional). */
	Delete(int, ...int) error
	/* Appends element to the end of the list. */
	Append(interface{}) error
	/* Pop and return the last element from the list. */
	Pop() (interface{}, error)
	/* Clear all elements from the list. */
	Clear() error

	/* Creates a copy of the current ListInterface. */
	Copy() (ListInterface, error)

	/* Returns a string representation of the list. */
	String() string

	/* Initializes the list. */
	Init()
}

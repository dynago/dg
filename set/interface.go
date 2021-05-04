package set

// SetInterface is the interface which defines whether a struct is a set or not.
type SetInterface interface {
	/* Return the number of elements in set. */
	Length() int
	/* Return the next key in set. */
	Iterate() <-chan interface{}

	/* Add element to the set. */
	Add(interface{}) error
	/* Remove element from the set. */
	Remove(interface{}) error
	/* Update the set, adding elements from the other set. */
	Combine(SetInterface) error
	/* Pop and return an arbitrary element from the set. */
	Pop() (interface{}, error)
	/* Clear all elements from the set. */
	Clear() error

	/* Test for membership in the set. */
	Contains(interface{}) (bool, error)
	/* Return true if the set has no elements in common with the other set. */
	Disjoint(SetInterface) (bool, error)
	/* Return true if the set has all elements in common with the other set. */
	Equals(SetInterface) (bool, error)
	/* Test whether every element in the other set is in the set. */
	SupersetOf(SetInterface) (bool, error)
	/* Test whether every element in the set is in the other set. */
	SubsetOf(SetInterface) (bool, error)

	/* Return a new set with elements common to the set and all others. */
	Intersection(SetInterface) (SetInterface, error)
	/* Return a new set with elements in either the set or the other but not both. */
	SymmetricDifference(SetInterface) (SetInterface, error)
	/* Return a new set with elements in the set that are not in the other set. */
	Difference(SetInterface) (SetInterface, error)
	/* Return a new set with elements from the set and the other set. */
	Union(SetInterface) (SetInterface, error)

	/* Creates a copy of the current SetInterface */
	Copy() (SetInterface, error)

	/* Return the value given a string representation of the bytes. */
	Get(string) interface{}

	/* Returns a string representation of the set. */
	String() string

	/* Initializes the set. */
	Init()
}

// Tuple implementation by Calder Lund

package tuple


type TupleInterface interface {
	/* Return the number of elements in tuple. */
	Length() int
	/* Return the next value in tuple. */
	Iterate() <-chan interface{}

	/* Test for membership in the tuple. */
	Contains(interface{}) (bool, error)
	/* Return true if the tuple has all elements in common with the other tuple. */
	Equals(TupleInterface) (bool, error)

	/* Return concatenation of two tuples together. */
	Concatenate(TupleInterface) (TupleInterface, error)
	/* Return tuple repeated n times. */
	Multiply(int) (TupleInterface, error)

	/* Returns the value at index. */
	Get(int) (interface{}, error)
	/* Returns the a tuple of values given range. */
	Range(int, int) (TupleInterface, error)
	/* Return first index of value. Returns -1 if not found. */
	Index(interface{}) (int, error)
	/* Return count of value. */
	Count(interface{}) (int, error)

	/* Creates a copy of the current TupleInterface. */
	Copy() (TupleInterface, error)

	/* Returns a string representation of the tuple. */
	String() string

	/* Initializes the tuple. */
	Init()
}

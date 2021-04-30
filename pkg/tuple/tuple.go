// Tuple implementation by Calder Lund

package tuple

import (
	"fmt"
	"strings"

	"internal/iterable"
)

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

type Tuple struct {
	values []interface{}
}

/* Return the number of elements in set. */
func (t *Tuple) Length() int {
	return len(t.values)
}

/* Return the next key in set. */
func (t *Tuple) Iterate() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for _, v := range t.values {
			c <- v
		}
		close(c)
	}()
	return c
}

/* Test for membership in the tuple. */
func (t *Tuple) Contains(value interface{}) (bool, error) {
	for _, v := range t.values {
		if v == value {
			return true, nil
		}
	}
	return false, nil
}

/* Return true if the tuple has all elements in common with the other tuple. */
func (t *Tuple) Equals(other TupleInterface) (bool, error) {
	if t.Length() != other.Length() {
		return false, nil
	}
	for i, tv := range t.values {
		ov, err := other.Get(i)
		if err != nil {
			return false, err
		}
		if tv != ov {
			return false, nil
		}
	}
	return true, nil
}

/* Return concatenation of two tuples together. */
func (t *Tuple) Concatenate(other TupleInterface) (TupleInterface, error) {
	output := new(Tuple)
	output.Init()
	output.values = t.values
	for value := range other.Iterate() {
		output.values = append(output.values, value)
	}
	return output, nil
}

/* Return tuple repeated n times. */
func (t *Tuple) Multiply(n int) (TupleInterface, error) {
	output := new(Tuple)
	output.Init()
	for i := 0; i < n; i++ {
		output.values = append(output.values, t.values...)
	}
	return output, nil
}

/* Returns the value at index. */
func (t *Tuple) Get(i int) (interface{}, error) {
	if i >= len(t.values) || i < 0 {
		return nil, fmt.Errorf("Input i out of range of tuple")
	}
	
	return t.values[i], nil
}

/* Returns a valid index given the length of the tuple. */
func validIndex(i int, length int) int {
	if i < 0 {
		return 0
	} 
	if i > length {
		return length
	}
	return i
}

/* Returns the a tuple of values given range. */
func (t *Tuple) Range(start int, end int) (TupleInterface, error) {
	start = validIndex(start, len(t.values))
	end = validIndex(end, len(t.values))

	output := new(Tuple)
	output.Init()
	for i := start; i < end; i++ {
		output.values = append(output.values, t.values[i])
	}
	return output, nil
}

/* Return first index of value. Returns -1 if not found. */
func (t *Tuple) Index(value interface{}) (int, error) {
	for i, v := range t.values {
		if v == value {
			return i, nil
		}
	}
	return -1, nil
}

/* Return count of value. */
func (t *Tuple) Count(value interface{}) (int, error) {
	count := 0
	for _, v := range t.values {
		if v == value {
			count += 1
		}
	}
	return count, nil
}

/* Creates a copy of the current TupleInterface. */
func (t *Tuple) Copy() (TupleInterface, error) {
	output, err := MakeTuple(t)
	return output, err
}

/* Returns a string representation of the tuple. */
func (t *Tuple) String() string {
	output := "("
	for _, value := range t.values {
		output += fmt.Sprintf("%v ", value)
	}
	output = strings.Trim(output, " ") + ")"
	return output
}

/* Initializes the tuple. */
func (t *Tuple) Init() {
	t.values = make([]interface{}, 0)
}

/* Initialize a new tuple object */
func MakeTuple(it ...iterable.Iterable) (TupleInterface, error) {
	output := new(Tuple)
	output.Init()
	if len(it) > 0 {
		i := it[0]
		for val := range i.Iterate() {
			output.values = append(output.values, val)
		}
	}
	return output, nil
}

/* Initialize a new tuple object */
func MakeTupleFromValues(values ...interface{}) (TupleInterface, error) {
	output := new(Tuple)
	output.Init()
	for _, val := range values {
		output.values = append(output.values, val)
	}
	return output, nil
}

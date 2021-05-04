// Package tuple implements tuples, which are collections which are ordered and immutable.
package tuple

import (
	"fmt"
	"strings"

	"github.com/CalderLund/DynamicGo/internal/helpers"
	"github.com/CalderLund/DynamicGo/internal/iterable"
)

// Tuple is a dynamic tuple structure.
type Tuple struct {
	values []interface{}
}

/* Length returns the number of elements in the tuple. */
func (t *Tuple) Length() int {
	return len(t.values)
}

/* Iterate returns the next key in the tuple. */
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

/* Contains tests for membership in the tuple. */
func (t *Tuple) Contains(value interface{}) (bool, error) {
	for _, v := range t.values {
		if v == value {
			return true, nil
		}
	}
	return false, nil
}

/* Equals returns true if the tuple has all elements in common with the other tuple. */
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

/* Concatenate returns concatenation of two tuples together. */
func (t *Tuple) Concatenate(other TupleInterface) (TupleInterface, error) {
	output := new(Tuple)
	output.Init()
	output.values = t.values
	for value := range other.Iterate() {
		output.values = append(output.values, value)
	}
	return output, nil
}

/* Multiply returns tuple repeated n times. */
func (t *Tuple) Multiply(n int) (TupleInterface, error) {
	output := new(Tuple)
	output.Init()
	for i := 0; i < n; i++ {
		output.values = append(output.values, t.values...)
	}
	return output, nil
}

/* Get returns the value at index. */
func (t *Tuple) Get(i int) (interface{}, error) {
	if i >= len(t.values) || i < 0 {
		return nil, fmt.Errorf("Index i out of range of tuple")
	}

	return t.values[i], nil
}

/* Range returns the a tuple of values given range. */
func (t *Tuple) Range(start int, end int) (TupleInterface, error) {
	start = helpers.ValidIndex(start, len(t.values))
	end = helpers.ValidIndex(end, len(t.values))

	output := new(Tuple)
	output.Init()
	for i := start; i < end; i++ {
		output.values = append(output.values, t.values[i])
	}
	return output, nil
}

/* Index returns first index of value. Returns -1 if not found. */
func (t *Tuple) Index(value interface{}) (int, error) {
	for i, v := range t.values {
		if v == value {
			return i, nil
		}
	}
	return -1, nil
}

/* Count returns count of value. */
func (t *Tuple) Count(value interface{}) (int, error) {
	count := 0
	for _, v := range t.values {
		if v == value {
			count += 1
		}
	}
	return count, nil
}

/* Copy creates a copy of the current TupleInterface. */
func (t *Tuple) Copy() (TupleInterface, error) {
	output, err := MakeTuple(t)
	return output, err
}

/* String returns a string representation of the tuple. */
func (t *Tuple) String() string {
	output := "("
	for _, value := range t.values {
		output += fmt.Sprintf("%v ", value)
	}
	output = strings.Trim(output, " ") + ")"
	return output
}

/* Init initializes the tuple. */
func (t *Tuple) Init() {
	t.values = make([]interface{}, 0)
}

/* MakeTuple initializes a new tuple object using an Iterable. */
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

/* MakeTupleFromValues initializes a new tuple object using interface{} objects. */
func MakeTupleFromValues(values ...interface{}) (TupleInterface, error) {
	output := new(Tuple)
	output.Init()
	for _, val := range values {
		output.values = append(output.values, val)
	}
	return output, nil
}

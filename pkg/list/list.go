// List implementation by Calder Lund

package list

import (
	"fmt"
	"strings"

	"internal/helpers"
	"internal/iterable"
)


type List struct {
	values []interface{}
}

/* Return the number of elements in list. */
func (l *List) Length() int {
	return len(l.values)
}

/* Return the next value in list. */
func (l *List) Iterate() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for _, v := range l.values {
			c <- v
		}
		close(c)
	}()
	return c
}

/* Test for membership in the list. */
func (l *List) Contains(value interface{}) (bool, error) {
	for _, v := range l.values {
		if v == value {
			return true, nil
		}
	}
	return false, nil
}

/* Return true if the list has all elements in common with the other list. */
func (l *List) Equals(other ListInterface) (bool, error) {
	if l.Length() != other.Length() {
		return false, nil
	}
	for i, lv := range l.values {
		ov, err := other.Get(i)
		if err != nil {
			return false, err
		}
		if lv != ov {
			return false, nil
		}
	}
	return true, nil
}

/* Return concatenation of two lists together. */
func (l *List) Concatenate(other ListInterface) (ListInterface, error) {
	output := new(List)
	output.Init()
	output.values = l.values
	for value := range other.Iterate() {
		output.values = append(output.values, value)
	}
	return output, nil
}

/* Return list repeated n times. */
func (l *List) Multiply(n int) (ListInterface, error) {
	output := new(List)
	output.Init()
	for i := 0; i < n; i++ {
		output.values = append(output.values, l.values...)
	}
	return output, nil
}

/* Return reversed list. */
func (l *List) Reverse() (ListInterface, error) {
	output := new(List)
	output.Init()
	n := len(l.values)
	for i := n - 1; i >= 0; i-- {
		output.values = append(output.values, l.values[i])
	}
	return output, nil
}

/* Returns the value at index. */
func (l *List) Get(i int) (interface{}, error) {
	if i >= len(l.values) || i < 0 {
		return nil, fmt.Errorf("Index i out of range of list")
	}
	
	return l.values[i], nil
}

/* Returns the a list of values given range. */
func (l *List) Range(start int, end int) (ListInterface, error) {
	start = helpers.ValidIndex(start, len(l.values))
	end = helpers.ValidIndex(end, len(l.values))

	output := new(List)
	output.Init()
	for i := start; i < end; i++ {
		output.values = append(output.values, l.values[i])
	}
	return output, nil
}

/* Return first index of value. Returns -1 if not found. */
func (l *List) Index(value interface{}) (int, error) {
	for i, v := range l.values {
		if v == value {
			return i, nil
		}
	}
	return -1, nil
}

/* Return count of value. */
func (l *List) Count(value interface{}) (int, error) {
	count := 0
	for _, v := range l.values {
		if v == value {
			count += 1
		}
	}
	return count, nil
}

/* Inserts the value at index. */
func (l *List) Insert(i int, value interface{}) error {
	if i >= len(l.values) || i < 0 {
		return fmt.Errorf("Index i out of range of list")
	}
	l.values[i] = value
	return nil
}

/* Sets values in given range to the values in the iterable. */
func (l *List) Set(start int, end int, it iterable.Iterable) error {
	if len(l.values) > 0 {
		s := helpers.ValidIndex(start, len(l.values))
		if s != start {
			return fmt.Errorf("Index start out of range of list")
		}
		e := helpers.ValidIndex(end, len(l.values))
		if e != end {
			return fmt.Errorf("Index end out of range of list")
		}
		
		c := it.Iterate()
		for i := start; i < end; i++ {
			val, ok := <-c
			if !ok {
				return fmt.Errorf("Iterable has no more values")
			} else {
				l.values[i] = val
			}
		}
		return nil
	}
	return fmt.Errorf("Cannot set on empty list")
}

/* Remove first occurence of the element from the list. */
func (l *List) Remove(value interface{}) error {
	i, err := l.Index(value)
	if err != nil {
		return err
	}
	if i >= 0 && i < len(l.values) - 1 {
		l.values = append(l.values[:i], l.values[i+1:]...)
	} else if i >= 0 {
		l.values = l.values[:i]
	}
	return nil
}

/* Removes range from the list. Takes two parameters: start (int), end (int: optional). */
func (l *List) Delete(start int, end ...int) error {
	var s, e int

	s = helpers.ValidIndex(start, len(l.values))
	if len(end) > 0 {
		e = helpers.ValidIndex(end[0], len(l.values))
	} else {
		e = helpers.ValidIndex(s + 1, len(l.values))
	}

	if e > s && e < len(l.values) {
		l.values = append(l.values[:s], l.values[e:]...)
	} else if e > s {
		l.values = l.values[:s]
	}
	
	return nil
}

/* Appends element to the end of the list. */
func (l *List) Append(value interface{}) error {
	l.values = append(l.values, value)
	return nil
}

/* Pop and return the last element from the list. */
func (l *List) Pop() (interface{}, error) {
	if len(l.values) > 0 {
		value := l.values[len(l.values)-1]
		l.values = l.values[:len(l.values)-1]
		return value, nil
	}
	return nil, fmt.Errorf("Cannot pop from empty list")
}

/* Clear all elements from the list. */
func (l *List) Clear() error {
	l.Init()
	return nil
}

/* Creates a copy of the current ListInterface. */
func (l *List) Copy() (ListInterface, error) {
	output, err := MakeList(l)
	return output, err
}

/* Returns a string representation of the list. */
func (l *List) String() string {
	output := "["
	for _, value := range l.values {
		output += fmt.Sprintf("%v ", value)
	}
	output = strings.Trim(output, " ") + "]"
	return output
}

/* Initializes the list. */
func (l *List) Init() {
	l.values = make([]interface{}, 0)
}

/* Initialize a new list object */
func MakeList(it ...iterable.Iterable) (ListInterface, error) {
	output := new(List)
	output.Init()
	if len(it) > 0 {
		i := it[0]
		for val := range i.Iterate() {
			output.values = append(output.values, val)
		}
	}
	return output, nil
}

/* Initialize a new list object */
func MakeListFromValues(values ...interface{}) (ListInterface, error) {
	output := new(List)
	output.Init()
	for _, val := range values {
		output.values = append(output.values, val)
	}
	return output, nil
}

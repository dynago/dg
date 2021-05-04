// Package set implements sets, which are collections which are both unordered and unindexed.
package set

import (
	"fmt"
	"strings"

	"github.com/CalderLund/DynamicGo/internal/helpers"
	"github.com/CalderLund/DynamicGo/internal/iterable"
)

// Set is a dynamic set structure.
type Set struct {
	values map[string]interface{}
}

/* Get returns the value given a string representation of the bytes. */
func (s *Set) Get(hash string) interface{} {
	value, ok := s.values[hash]
	if !ok {
		return nil
	}
	return value
}

/* Length returns the number of elements in set. */
func (s *Set) Length() int {
	return len(s.values)
}

/* Iterate returns the next key in set. */
func (s *Set) Iterate() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for _, v := range s.values {
			c <- v
		}
		close(c)
	}()
	return c
}

/* Add adds element to the set. */
func (s *Set) Add(value interface{}) error {
	hash, err := helpers.GetSHA(value)
	if err != nil {
		return err
	}
	s.values[hash] = value
	return nil
}

/* Remove removes element from the set. */
func (s *Set) Remove(value interface{}) error {
	hash, err := helpers.GetSHA(value)
	if err != nil {
		return err
	}
	delete(s.values, hash)
	return nil
}

/* Combine updates the set, adding elements from the other set. */
func (s *Set) Combine(other SetInterface) error {
	for value := range other.Iterate() {
		hash, err := helpers.GetSHA(value)
		if err != nil {
			return err
		}
		s.values[hash] = value
	}
	return nil
}

/* Pop pops and return an arbitrary element from the set. */
func (s *Set) Pop() (interface{}, error) {
	for value := range s.Iterate() {
		hash, err := helpers.GetSHA(value)
		if err != nil {
			return nil, err
		}
		delete(s.values, hash)
		return value, nil
	}
	return nil, nil
}

/* Clear clears all elements from the set. */
func (s *Set) Clear() error {
	s.Init()
	return nil
}

/* Contains tests for membership in the set. */
func (s *Set) Contains(value interface{}) (bool, error) {
	hash, err := helpers.GetSHA(value)
	if err != nil {
		return false, err
	}
	_, ok := s.values[hash]
	return ok, nil
}

/* Disjoint returns true if the set has no elements in common with the other set. */
func (s *Set) Disjoint(other SetInterface) (bool, error) {
	for value := range other.Iterate() {
		hash, err := helpers.GetSHA(value)
		if err != nil {
			return false, err
		}
		if _, ok := s.values[hash]; ok {
			return false, nil
		}
	}
	return true, nil
}

/* Equals returns true if the set has all elements in common with the other set. */
func (s *Set) Equals(other SetInterface) (bool, error) {
	var ok1, ok2 bool
	var err error
	if ok1, err = other.SupersetOf(s); err != nil {
		return false, err
	}
	if ok2, err = s.SupersetOf(other); err != nil {
		return false, err
	}
	return ok1 && ok2, nil
}

/* SupersetOf tests whether every element in the other set is in the set. */
func (s *Set) SupersetOf(other SetInterface) (bool, error) {
	for value := range other.Iterate() {
		hash, err := helpers.GetSHA(value)
		if err != nil {
			return false, err
		}
		if _, ok := s.values[hash]; !ok {
			return false, nil
		}
	}
	return true, nil
}

/* SubsetOf tests whether every element in the set is in the other set. */
func (s *Set) SubsetOf(other SetInterface) (bool, error) {
	ok, err := other.SupersetOf(s)
	return ok, err
}

/* Intersection returns a new set with elements common to the set and all others. */
func (s *Set) Intersection(other SetInterface) (SetInterface, error) {
	output, err := MakeSet()
	if err != nil {
		return nil, err
	}
	for value := range other.Iterate() {
		hash, err := helpers.GetSHA(value)
		if err != nil {
			return nil, err
		}
		if v, ok := s.values[hash]; ok {
			if err = output.Add(v); err != nil {
				return nil, err
			}
		}
	}
	return output, nil
}

/* SymmetricDifference returns a new set with elements in either the set or the other but not both. */
func (s *Set) SymmetricDifference(other SetInterface) (SetInterface, error) {
	output, err1 := s.Difference(other)
	if err1 != nil {
		return nil, err1
	}
	otherOutput, err2 := other.Difference(s)
	if err2 != nil {
		return nil, err2
	}
	err := output.Combine(otherOutput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

/* Difference returns a new set with elements in the set that are not in the other set. */
func (s *Set) Difference(other SetInterface) (SetInterface, error) {
	output, err := s.Copy()
	if err != nil {
		return nil, err
	}
	for value := range other.Iterate() {
		hash, err := helpers.GetSHA(value)
		if err != nil {
			return nil, err
		}
		if v, ok := s.values[hash]; ok {
			err := output.Remove(v)
			if err != nil {
				return nil, err
			}
		}
	}
	return output, nil
}

/* Union returns a new set with elements from the set and the other set. */
func (s *Set) Union(other SetInterface) (SetInterface, error) {
	output, err := s.Copy()
	if err != nil {
		return nil, err
	}
	output.Combine(other)
	return output, nil
}

/* Copy creates a copy of the current SetInterface */
func (s *Set) Copy() (SetInterface, error) {
	output, err := MakeSet(s)
	return output, err
}

/* String returns a string representation of the set. */
func (s *Set) String() string {
	output := "("
	for _, value := range s.values {
		output += fmt.Sprintf("%v ", value)
	}
	output = strings.Trim(output, " ") + ")"
	return output
}

/* Init initializes the set. */
func (s *Set) Init() {
	s.values = make(map[string]interface{})
}

/* MakeSet initializes a new set object using an Iterable. */
func MakeSet(it ...iterable.Iterable) (SetInterface, error) {
	output := new(Set)
	output.Init()
	if len(it) > 0 {
		i := it[0]
		for val := range i.Iterate() {
			if err := output.Add(val); err != nil {
				return nil, err
			}
		}
	}
	return output, nil
}

/* MakeSetFromValues initializes a new set object using interface{} objects. */
func MakeSetFromValues(values ...interface{}) (SetInterface, error) {
	output := new(Set)
	output.Init()
	for _, val := range values {
		if err := output.Add(val); err != nil {
			return nil, err
		}
	}
	return output, nil
}

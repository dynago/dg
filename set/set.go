// Set implementation by Calder Lund

package set

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type SetInterface interface {
	/* Return the number of elements in set. */
	Length() int
	/* Return the next key in set. */
	Iterate() <-chan string

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

type Set struct {
	values map[string]interface{}
}

func getSHA(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	hasher := sha1.New()
	hasher.Write(b)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

/* Return the value given a string representation of the bytes. */
func (s *Set) Get(hash string) interface{} {
	value, ok := s.values[hash]
	if !ok {
		return nil
	}
	return value
}

/* Return the number of elements in set. */
func (s *Set) Length() int {
	return len(s.values)
}

/* Return the next key in set. */
func (s *Set) Iterate() <-chan string {
	c := make(chan string)
	go func() {
		for k := range s.values {
			c <- k
		}
		close(c)
	}()
	return c
}

/* Add element to the set. */
func (s *Set) Add(value interface{}) error {
	hash, err := getSHA(value)
	if err != nil {
		return err
	}
	s.values[hash] = value
	return nil
}

/* Remove element from the set. */
func (s *Set) Remove(value interface{}) error {
	hash, err := getSHA(value)
	if err != nil {
		return err
	}
	delete(s.values, hash)
	return nil
}

/* Update the set, adding elements from the other set. */
func (s *Set) Combine(other SetInterface) error {
	for hash := range other.Iterate() {
		value := other.Get(hash)
		s.values[hash] = value
	}
	return nil
}

/* Pop and return an arbitrary element from the set. */
func (s *Set) Pop() (interface{}, error) {
	for hash := range s.Iterate() {
		value := s.values[hash]
		delete(s.values, hash)
		return value, nil
	}
	return nil, nil
}

/* Clear all elements from the set. */
func (s *Set) Clear() error {
	for hash := range s.Iterate() {
		delete(s.values, hash)
	}
	return nil
}

/* Test for membership in the set. */
func (s *Set) Contains(value interface{}) (bool, error) {
	hash, err := getSHA(value)
	if err != nil {
		return false, err
	}
	_, ok := s.values[hash]
	return ok, nil
}

/* Return true if the set has no elements in common with the other set. */
func (s *Set) Disjoint(other SetInterface) (bool, error) {
	for hash := range other.Iterate() {
		if _, ok := s.values[hash]; ok {
			return false, nil
		}
	}
	return true, nil
}

/* Return true if the set has all elements in common with the other set. */
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

/* Test whether every element in the other set is in the set. */
func (s *Set) SupersetOf(other SetInterface) (bool, error) {
	for hash := range other.Iterate() {
		if _, ok := s.values[hash]; !ok {
			return false, nil
		}
	}
	return true, nil
}

/* Test whether every element in the set is in the other set. */
func (s *Set) SubsetOf(other SetInterface) (bool, error) {
	ok, err := other.SupersetOf(s)
	return ok, err
}

/* Return a new set with elements common to the set and all others. */
func (s *Set) Intersection(other SetInterface) (SetInterface, error) {
	output := MakeSet()
	for hash := range other.Iterate() {
		if value, ok := s.values[hash]; ok {
			err := output.Add(value)
			if err != nil {
				return nil, err
			}
		}
	}
	return output, nil
}

/* Return a new set with elements in either the set or the other but not both. */
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

/* Return a new set with elements in the set that are not in the other set. */
func (s *Set) Difference(other SetInterface) (SetInterface, error) {
	output, err := s.Copy()
	if err != nil {
		return nil, err
	}
	for hash := range other.Iterate() {
		if value, ok := s.values[hash]; ok {
			err := output.Remove(value)
			if err != nil {
				return nil, err
			}
		}
	}
	return output, nil
}

/* Return a new set with elements from the set and the other set. */
func (s *Set) Union(other SetInterface) (SetInterface, error) {
	output, err := s.Copy()
	if err != nil {
		return nil, err
	}
	output.Combine(other)
	return output, nil
}

/* Creates a copy of the current SetInterface */
func (s *Set) Copy() (SetInterface, error) {
	output := MakeSet()
	for hash := range s.Iterate() {
		err := output.Add(s.values[hash])
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

/* Returns a string representation of the set. */
func (s *Set) String() string {
	output := "("
	for _, value := range s.values {
		output += fmt.Sprintf("%v ", value)
	}
	output = strings.Trim(output, " ") + ")"
	return output
}

/* Initializes the set. */
func (s *Set) Init() {
	s.values = make(map[string]interface{})
}

/* Initialize a new set object */
func MakeSet() SetInterface {
	output := new(Set)
	output.Init()
	return output
}

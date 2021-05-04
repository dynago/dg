// Package dict implements dictionaries, which are collections of unordered key/value pairs.
package dict

import (
	"fmt"
	"strings"

	"github.com/CalderLund/DynamicGo/internal/helpers"
	"github.com/CalderLund/DynamicGo/internal/iterable"
	"github.com/CalderLund/DynamicGo/list"
	"github.com/CalderLund/DynamicGo/tuple"
)

// Dict is a dynamic dictionary structure.
type Dict struct {
	keys   map[string]interface{} // hash of key to key
	values map[string]interface{} // hash of key to value
}

/* Length returns the number of elements in dict. */
func (d *Dict) Length() int {
	return len(d.keys)
}

/* Iterate returns the next key in dict. */
func (d *Dict) Iterate() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for _, v := range d.keys {
			c <- v
		}
		close(c)
	}()
	return c
}

/* Remove removes element from the dict. */
func (d *Dict) Remove(key interface{}) error {
	hash, err := helpers.GetSHA(key)
	if err != nil {
		return err
	}
	delete(d.keys, hash)
	delete(d.values, hash)
	return nil
}

/* Get returns the value with given key. */
func (d *Dict) Get(key interface{}) (interface{}, error) {
	hash, err := helpers.GetSHA(key)
	if err != nil {
		return nil, err
	}
	value, ok := d.values[hash]
	if !ok {
		return nil, nil
	}
	return value, nil
}

/* Set sets the value at given key to given value. */
func (d *Dict) Set(key interface{}, value interface{}) error {
	hash, err := helpers.GetSHA(key)
	if err != nil {
		return err
	}
	d.keys[hash] = key
	d.values[hash] = value
	return nil
}

/* Combine updates the dict, adding elements from the other dict. Old values are replaced with new. */
func (d *Dict) Combine(other DictInterface) error {
	for key := range other.Iterate() {
		value, err := other.Get(key)
		if err != nil {
			return err
		}
		if err = d.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

/* PopKey pops and returns an arbitrary key from the dict. */
func (d *Dict) PopKey() (interface{}, error) {
	for key := range d.Iterate() {
		if err := d.Remove(key); err != nil {
			return nil, err
		}
		return key, nil
	}
	return nil, nil
}

/* PopValue pops and returns an arbitrary value from the dict. */
func (d *Dict) PopValue() (interface{}, error) {
	for key := range d.Iterate() {
		value, err := d.Get(key)
		if err != nil {
			return nil, err
		}
		if err = d.Remove(key); err != nil {
			return nil, err
		}
		return value, nil
	}
	return nil, nil
}

/* Pop pops and returns an arbitrary item from the dict. */
func (d *Dict) Pop() (interface{}, interface{}, error) {
	for key := range d.Iterate() {
		value, err := d.Get(key)
		if err != nil {
			return nil, nil, err
		}
		if err = d.Remove(key); err != nil {
			return nil, nil, err
		}
		return key, value, nil
	}
	return nil, nil, nil
}

/* Clear clears all elements from the dict. */
func (d *Dict) Clear() error {
	d.Init()
	return nil
}

/* Contains tests for membership in the dict. */
func (d *Dict) Contains(key interface{}) (bool, error) {
	hash, err := helpers.GetSHA(key)
	if err != nil {
		return false, err
	}
	_, ok := d.keys[hash]
	return ok, nil
}

/* Equals returns true if the dict has all elements in common with the other dict. */
func (d *Dict) Equals(other DictInterface) (bool, error) {
	if d.Length() != other.Length() {
		return false, nil
	}
	for ok := range other.Iterate() {
		ov, err1 := other.Get(ok)
		if err1 != nil {
			return false, err1
		}
		dv, err2 := d.Get(ok)
		if err2 != nil {
			return false, err2
		}
		if dv != ov {
			return false, nil
		}
	}
	return true, nil
}

/* Keys returns a tuple of keys. */
func (d *Dict) Keys() (tuple.TupleInterface, error) {
	keys := make([]interface{}, 0)
	for _, key := range d.keys {
		keys = append(keys, key)
	}
	tup, err := tuple.MakeTupleFromValues(keys...)
	if err != nil {
		return nil, err
	}
	return tup, nil
}

/* Values returns a tuple of values. */
func (d *Dict) Values() (tuple.TupleInterface, error) {
	values := make([]interface{}, 0)
	for _, value := range d.values {
		values = append(values, value)
	}
	tup, err := tuple.MakeTupleFromValues(values...)
	if err != nil {
		return nil, err
	}
	return tup, nil
}

/* Items returns a tuple of key/value pairs. */
func (d *Dict) Items() (tuple.TupleInterface, error) {
	l, err := list.MakeList()
	if err != nil {
		return nil, err
	}
	var value interface{}
	for _, key := range d.keys {
		value, err = d.Get(key)
		if err != nil {
			return nil, err
		}
		tup, errt := tuple.MakeTupleFromValues(key, value)
		if errt != nil {
			return nil, errt
		}
		if err = l.Append(tup); err != nil {
			return nil, err
		}
	}
	t, errt := tuple.MakeTuple(l)
	return t, errt
}

/* Copy creates a copy of the current DictInterface */
func (d *Dict) Copy() (DictInterface, error) {
	keys := make([]interface{}, 0)
	values := make([]interface{}, 0)
	for _, key := range d.keys {
		value, err := d.Get(key)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
		values = append(values, value)
	}
	output, err := MakeDictFromKeyValues(keys, values)
	return output, err
}

/* String returns a string representation of the dict. */
func (d *Dict) String() string {
	output := "{"
	for key := range d.Iterate() {
		value, err := d.Get(key)
		if err != nil {
			return ""
		}
		output += fmt.Sprintf("(%v %v) ", key, value)
	}
	output = strings.Trim(output, " ") + "}"
	return output
}

/* Init initializes the dict. */
func (d *Dict) Init() {
	d.keys = make(map[string]interface{})
	d.values = make(map[string]interface{})
}

/* MakeDict initializes a new dict object using an Iterable. Every even-indexed element is a key and odd-indexed element is a value. */
func MakeDict(it ...iterable.Iterable) (DictInterface, error) {
	output := new(Dict)
	output.Init()
	if len(it) > 0 {
		i := 0
		iter := it[0]
		var key interface{}
		for v := range iter.Iterate() {
			if i%2 == 0 {
				key = v
			} else {
				output.Set(key, v)
			}
			i += 1
		}
	}
	return output, nil
}

/* MakeDictFromKeyValues initializes a new dict object using keys and values. */
func MakeDictFromKeyValues(keys []interface{}, values []interface{}) (DictInterface, error) {
	output := new(Dict)
	output.Init()
	if len(keys) != len(values) {
		return nil, fmt.Errorf("Number of keys does not match number of value")
	}
	for i, key := range keys {
		value := values[i]
		output.Set(key, value)
	}
	return output, nil
}

/* MakeDictFromItems initializes a new dict object using . */
func MakeDictFromItems(items ...tuple.TupleInterface) (DictInterface, error) {
	output := new(Dict)
	output.Init()
	for _, item := range items {
		if item.Length() != 2 {
			return nil, fmt.Errorf("Each item must be of length 2 (key, value)")
		}
		key, err0 := item.Get(0)
		if err0 != nil {
			return nil, err0
		}
		value, err1 := item.Get(1)
		if err1 != nil {
			return nil, err1
		}
		output.Set(key, value)
	}
	return output, nil
}
